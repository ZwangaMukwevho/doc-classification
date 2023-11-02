package model

type Attachment struct {
	ID         string
	Name       string
	Bytestream string
	MimeType   string
}

type Message struct {
	ID        string
	Subject   string
	Timestamp string
	File      Attachment
}
