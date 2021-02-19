package main

import (
	"fmt"
	"log"

	"github.com/ldej/go-acapy-client"
)

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

func (app *App) CredentialExchangeEventHandler(event acapy.CredentialExchangeRecord) {
	connection, _ := app.acapy.GetConnection(event.ConnectionID)
	fmt.Printf("\n -> Credential Exchange update: %s - %s - %s\n", event.CredentialExchangeID, connection.TheirLabel, event.State)
}

func (app *App) RevocationRegistryEventHandler(event acapy.RevocationRegistry) {
	fmt.Printf("\n -> Revocation Registry update: %s - %s\n", event.RevocationRegistryID, event.State)
}

func (app *App) ProblemReportEventHandler(event acapy.ProblemReportEvent) {
	fmt.Printf("\n -> Received problem report: %+v\n", event)
}

func (app *App) PresentationExchangeEventHandler(event acapy.PresentationExchangeRecord) {
	connection, _ := app.acapy.GetConnection(event.ConnectionID)
	fmt.Printf("\n -> Presentation Exchange update: %s - %s - %s\n", connection.TheirLabel, event.PresentationExchangeID, event.State)
}

func (app *App) CredentialRevocationEventHandler(event acapy.CredentialRevocationRecord) {
	fmt.Printf("\n -> Issuer Credential Revocation: %s - %s - %s\n", event.CredentialExchangeID, event.RecordID, event.State)
}

func (app *App) CredentialExchangeV2EventHandler(event acapy.CredentialExchangeRecordV2) {
	connection, _ := app.acapy.GetConnection(event.ConnectionID)
	fmt.Printf("\n -> Credential Exchange V2 update: %s - %s - %s\n", event.CredentialExchangeID, connection.TheirLabel, event.State)
}

func (app *App) CredentialExchangeDIFEventHandler(event acapy.CredentialExchangeDIF) {
	record, _ := app.acapy.GetCredentialExchangeV2(event.CredentialExchangeID)
	connection, _ := app.acapy.GetConnection(record.CredentialExchangeRecord.ConnectionID)
	fmt.Printf("\n -> Credential Exchange DIF Event: %s - %s - %s", connection.TheirLabel, event.CredentialExchangeID, event.State)
}

func (app *App) CredentialExchangeIndyEventHandler(event acapy.CredentialExchangeIndy) {
	record, _ := app.acapy.GetCredentialExchangeV2(event.CredentialExchangeID)
	connection, _ := app.acapy.GetConnection(record.CredentialExchangeRecord.ConnectionID)
	fmt.Printf("\n -> Credential Exchange Indy Event: %s - %s - %s", connection.TheirLabel, event.CredentialExchangeID, event.CredentialExchangeIndyID)
}

func (app *App) BasicMessagesEventHandler(event acapy.BasicMessagesEvent) {
	connection, _ := app.acapy.GetConnection(event.ConnectionID)
	fmt.Printf("\n -> Received message from %q (%s): %s\n", connection.TheirLabel, event.ConnectionID, event.Content)
}

func (app *App) PingEventHandler(event acapy.PingEvent) {
	fmt.Printf("\n -> Ping Event: %q state: %q responded: %t\n", event.ConnectionID, event.State, event.Responded)
}
