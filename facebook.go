package fsmfacebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-carrot/fsm"
	"github.com/go-carrot/fsm-dynamo-store"
	"github.com/julienschmidt/httprouter"
	"github.com/tylerb/graceful"
)

func Start(stateMachine fsm.StateMachine, startState string) {
	// Create Store
	store := &dynamostore.DynamoStore{
		Network:     "facebook",
		DynamoTable: os.Getenv("DYNAMO_DB_TABLE"),
		DynamoSession: dynamodb.New(
			session.New(
				&aws.Config{
					Region: aws.String(os.Getenv("DYNAMO_DB_REGION")),
				},
			),
		),
	}

	// Build Server
	srv := &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:    ":" + os.Getenv("PORT"),
			Handler: buildRouter(store, stateMachine, startState),
		},
	}
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func buildRouter(store fsm.Store, stateMachine fsm.StateMachine, startState string) *httprouter.Router {
	// Router
	router := &httprouter.Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
	}
	router.HandlerFunc(http.MethodGet, "/facebook", FacebookSetupWebhook)
	router.HandlerFunc(http.MethodPost, "/facebook", GetFacebookWebhook(store, stateMachine, startState))
	return router
}

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

func GetFacebookWebhook(store fsm.Store, stateMachine fsm.StateMachine, startState string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Body
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

				// Get Traverser
				newTraverser := false
				traverser, err := store.FetchTraverser(messagingEvent.Sender.ID)
				if err != nil {
					traverser, _ = store.CreateTraverser(messagingEvent.Sender.ID)
					traverser.SetCurrentState(startState)
					newTraverser = true
				}

				// Create Emitter
				emitter := &FacebookEmitter{
					UUID: traverser.UUID(),
				}

				// Get Current State
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
