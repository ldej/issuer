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
	issuerPort := os.Getenv("CONTROLLER_PORT")

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
		{
			var v1 = api.PathPrefix("/v1").Subrouter()
			v1.HandleFunc("/create-invitation", app.createInvitationV1).Methods(http.MethodPost)
			v1.HandleFunc("/schema", app.registerSchema).Methods(http.MethodPost)
			v1.HandleFunc("/credential-definition", app.createCredentialDefinition).Methods(http.MethodPost)
			v1.HandleFunc("/issue-credential", app.issueCredentialV1).Methods(http.MethodPost)
		}
		{
			var v2 = api.PathPrefix("/v2").Subrouter()
			v2.HandleFunc("/create-invitation", app.createInvitationV2).Methods(http.MethodPost)
			v2.HandleFunc("/issue-credential", app.issueCredentialV2).Methods(http.MethodPost)
		}
	}

	r.HandleFunc("/webhooks/topic/{topic}/", acapy.CreateWebhooksHandler(acapy.WebhookHandlers{
		ConnectionsEventHandler:            app.ConnectionsEventHandler,
		BasicMessagesEventHandler:          app.BasicMessagesEventHandler,
		ProblemReportEventHandler:          app.ProblemReportEventHandler,
		CredentialExchangeEventHandler:     app.CredentialExchangeEventHandler,
		CredentialExchangeV2EventHandler:   app.CredentialExchangeV2EventHandler,
		CredentialExchangeDIFEventHandler:  app.CredentialExchangeDIFEventHandler,
		CredentialExchangeIndyEventHandler: app.CredentialExchangeIndyEventHandler,
		RevocationRegistryEventHandler:     app.RevocationRegistryEventHandler,
		PresentationExchangeEventHandler:   app.PresentationExchangeEventHandler,
		CredentialRevocationEventHandler:   app.CredentialRevocationEventHandler,
		PingEventHandler:                   app.PingEventHandler,
		OutOfBandEventHandler:              app.OutOfBandEventHandler,
	}))
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

func (app *App) createInvitationV1(w http.ResponseWriter, r *http.Request) {
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

func (app *App) createInvitationV2(w http.ResponseWriter, r *http.Request) {
	invitation, err := app.acapy.CreateOutOfBandInvitation(
		acapy.CreateOutOfBandInvitationRequest{
			HandshakeProtocols: acapy.DefaultHandshakeProtocols,
		},
		true,
		false,
	)
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
	var request = struct {
		Name       string   `json:"name"`
		Version    string   `json:"version"`
		Attributes []string `json:"attributes"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Failed to decode request: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	schema, err := app.acapy.RegisterSchema(
		request.Name,
		request.Version,
		request.Attributes,
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
		SchemaID               string `json:"schema_id"`
		SupportRevocation      bool   `json:"support_revocation"`
		RevocationRegistrySize int    `json:"revocation_registry_size"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Failed to decode request: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	credentialDefinitionID, err := app.acapy.CreateCredentialDefinition("", request.SupportRevocation, request.RevocationRegistrySize, request.SchemaID)
	if err != nil {
		log.Printf("Failed to create credential definition: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(credentialDefinitionID))
}

func (app *App) issueCredentialV1(w http.ResponseWriter, r *http.Request) {
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
		request.ConnectionID,
		acapy.NewCredentialPreview(attributes),
		request.Comment,
		request.CredentialDefinitionID,
		request.IssuerDID,
		request.SchemaID,
	)
	if err != nil {
		log.Printf("Failed to issue credential: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (app *App) issueCredentialV2(w http.ResponseWriter, r *http.Request) {
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
	var attributes []acapy.CredentialPreviewAttributeV2
	for key, value := range request.Attributes {
		attributes = append(attributes, acapy.CredentialPreviewAttributeV2{
			Name:  key,
			Value: value,
		})
	}

	_, err = app.acapy.IssueCredentialV2(
		request.ConnectionID,
		acapy.NewCredentialPreviewV2(attributes),
		request.Comment,
		request.CredentialDefinitionID,
		request.IssuerDID,
		request.SchemaID,
	)
	if err != nil {
		log.Printf("Failed to issue credential: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
