package models

type Message struct {
	Text  string
	Whose string
}

type Party struct {
	Name      string
	City      string
	Organizer string
	Date      string
	About     string
	Id        string
	Messages  []Message
}

type Events struct {
	Parties []Party
}
