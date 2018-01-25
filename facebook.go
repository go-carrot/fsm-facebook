package fsmfacebook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-carrot/fsm"
)

// GetFacebookSetupWebhook adds support for the Messenger Platform's webhook verification
// to your webhook. This is required to ensure your webhook is authentic and working.
//
//  This must be a GET request, and have the same URL as the POST request.
//
// https://developers.facebook.com/docs/messenger-platform/getting-started/webhook-setup
func FacebookSetupWebhook(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	mode := queryParams.Get("hub.mode")
	challenge := queryParams.Get("hub.challenge")
	verifyToken := queryParams.Get("hub.verify_token")

	if mode == "subscribe" && verifyToken == os.Getenv("FACEBOOK_VERIFY_TOKEN") {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

// GetFacebookWebhook is the webhook that facebook posts to when a message is
// received from a user.
//
// This must be a POST request, and have the same URL as the GET request.
//
// https://developers.facebook.com/docs/messenger-platform/getting-started/webhook-setup
func GetFacebookWebhook(store fsm.Store, stateMachine fsm.StateMachine) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get body
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()

		// Parse body into struct
		cb := new(MessageReceivedCallback)
		json.Unmarshal([]byte(body), cb)

		// For each entry
		for _, i := range cb.Entry {
			// Iterate over each messaging event
			for _, messagingEvent := range i.MessagingEvents {

				// Get traverser
				newTraverser := false
				traverser, err := store.FetchTraverser(messagingEvent.Sender.ID)
				if err != nil {
					traverser, _ = store.CreateTraverser(messagingEvent.Sender.ID)
					traverser.SetCurrentState("start")
					newTraverser = true
				}

				// Create emitter
				emitter := &FacebookEmitter{
					UUID: traverser.UUID(),
				}

				// Get current state
				currentState := stateMachine[traverser.CurrentState()](emitter, traverser)
				if newTraverser {
					currentState.EntryAction()
				}

				// Transition
				newState := currentState.Transition(messagingEvent.Message.Text)
				err = newState.EntryAction()
				if err == nil {
					traverser.SetCurrentState(newState.Slug)
				}
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
