package fsmfacebook

import "github.com/BrandonRomano/wrecker"
import "os"

type MessageData struct {
	Recipient MessageRecipient `json:"recipient"`
	Message   SendMessageData  `json:"message"`
}

type MessageRecipient struct {
	ID string `json:"id"`
}

type SendMessageData struct {
	Text string `json:"text"`
}

type FacebookEmitter struct {
	UUID string
}

func (f *FacebookEmitter) Emit(input interface{}) error {
	MessageData := MessageData{
		Recipient: MessageRecipient{
			ID: f.UUID,
		},
		Message: SendMessageData{
			Text: input.(string),
		},
	}

	client := wrecker.New("https://graph.facebook.com/v2.6")
	client.Post("/me/messages").
		URLParam("access_token", os.Getenv("FACEBOOK_ACCESS_TOKEN")).
		Body(MessageData).
		Execute()

	return nil
}
