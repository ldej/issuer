package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
		acapy: acapy.NewClient(fmt.Sprintf("http://acapy:%s", acapyAdminPort), ""),
	}

	r := mux.NewRouter()
	{
		api := r.PathPrefix("/api").Subrouter()
		api.HandleFunc("/create-invitation", app.createInvitation)
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

	server := &http.Server{
		Addr:    ":" + issuerPort,
		Handler: r,
	}

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
		w.Write([]byte(err.Error()))
		return
	}

	png, err := qrcode.Encode(invitation.InvitationURL, qrcode.Medium, 256)
	if err != nil {
		log.Printf("Failed to create qr code: %s", err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(png)
	w.Header().Set("Content-Type", "image/png")
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
