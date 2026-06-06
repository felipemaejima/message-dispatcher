package domain

type Message struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

type Notification struct {
	Channel string `json:"channel"`
	Message Message
}
