package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ldej/go-acapy-client"
	"github.com/skip2/go-qrcode"
)

type App struct {
	acapy *acapy.Client
}

func main() {
	acapyAdminPort := os.Getenv("ACAPY_ADMIN_PORT")
	issuerPort := os.Getenv("ISSUER_PORT")

	app := App{
		acapy: acapy.NewClient(fmt.Sprintf("http://acapy:%s", acapyAdminPort)),
	}

	ready, err := app.acapy.IsReady()
	if err != nil {
		log.Fatalf("Error while checking if ACA-py is ready: %s", err.Error())
		return
	} else if !ready {
		log.Fatalf("ACA-py has started but it not ready")
		return
	} else {
		log.Println("ACA-py is ready on port", acapyAdminPort)
	}

	r := mux.NewRouter()
	{
		api := r.PathPrefix("/api").Subrouter()
		api.HandleFunc("/create-invitation", app.createInvitation).Methods(http.MethodPost)
		api.HandleFunc("/schema", app.registerSchema).Methods(http.MethodPost)
		api.HandleFunc("/credential-definition", app.createCredentialDefinition).Methods(http.MethodPost)
		api.HandleFunc("/issue-credential", app.issueCredential).Methods(http.MethodPost)
	}

	r.HandleFunc("/webhooks/topic/{topic}/", acapy.WebhookHandler(
		app.ConnectionsEventHandler,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		app.OutOfBandEventHandler,
	))
	r.NotFoundHandler = http.HandlerFunc(NotFound)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	server := &http.Server{
		Addr:    ":" + issuerPort,
		Handler: loggedRouter,
	}

	log.Println("Listening on port", issuerPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not Found: %s", r.URL.String())
}

func (app *App) createInvitation(w http.ResponseWriter, r *http.Request) {
	invitation, err := app.acapy.CreateInvitation("", true, false, false)
	if err != nil {
		log.Printf("Failed to create invitation: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	png, err := qrcode.Encode(invitation.InvitationURL, qrcode.Medium, 256)
	if err != nil {
		log.Printf("Failed to create qr code: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(png)
	w.Header().Set("Content-Type", "image/png")
}

func (app *App) registerSchema(w http.ResponseWriter, r *http.Request) {
	schema, err := app.acapy.RegisterSchema(
		"ldej",
		"1.0",
		[]string{"date"},
	)
	if err != nil {
		log.Printf("Failed to register schema: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(schema)
}

func (app *App) createCredentialDefinition(w http.ResponseWriter, r *http.Request) {
	var request = struct {
		SchemaID string `json:"schema_id"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Failed to decode request: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	credentialDefinitionID, err := app.acapy.CreateCredentialDefinition("", true, 4, request.SchemaID)
	if err != nil {
		log.Printf("Failed to create credential definition: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(credentialDefinitionID))
}

func (app *App) issueCredential(w http.ResponseWriter, r *http.Request) {
	var request = struct {
		ConnectionID           string            `json:"connection_id"`
		CredentialDefinitionID string            `json:"credential_definition_id"`
		Attributes             map[string]string `json:"attributes"`
		Comment                string            `json:"comment"`
		IssuerDID              string            `json:"issuer_did"`
		SchemaID               string            `json:"schema_id"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Failed to decode request: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var attributes []acapy.CredentialPreviewAttribute
	for key, value := range request.Attributes {
		attributes = append(attributes, acapy.CredentialPreviewAttribute{
			Name:  key,
			Value: value,
		})
	}

	_, err = app.acapy.IssueCredential(
		request.CredentialDefinitionID,
		request.ConnectionID,
		"", //request.IssuerDID,
		request.Comment,
		acapy.NewCredentialPreview(attributes),
		"", //request.SchemaID,
	)
	if err != nil {
		log.Printf("Failed to issue credential: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (app *App) ConnectionsEventHandler(event acapy.Connection) {
	if event.Alias == "" {
		connection, _ := app.acapy.GetConnection(event.ConnectionID)
		event.Alias = connection.TheirLabel
	}
	log.Printf(" -> Connection %q (%s), update to state %q rfc23 state %q", event.Alias, event.ConnectionID, event.State, event.RFC23State)
}

func (app *App) OutOfBandEventHandler(event acapy.OutOfBandEvent) {
	log.Printf(" -> Out of Band Event: %q state %q", event.InvitationID, event.State)
}
