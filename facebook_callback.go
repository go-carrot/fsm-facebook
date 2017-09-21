package fsmfacebook

type MessageReceivedCallback struct {
	Object string         `json:"object"`
	Entry  []MessageEntry `json:"entry"`
}

type MessageEntry struct {
	ID              string           `json:"id"`
	Time            int64            `json:"time"`
	MessagingEvents []MessagingEvent `json:"messaging"`
}

type MessagingEvent struct {
	Sender    Sender    `json:"sender"`
	Recipient Recipient `json:"recipient"`
	Timestamp int64     `json:"timestamp"`
	Message   Message   `json:"message"`
}

type Sender struct {
	ID string `json:"id"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	MessageID string `json:"mid"`
	Sequence  int64  `json:"seq"`
	Text      string `json:"text"`
}
