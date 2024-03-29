package model

type Attachment struct {
	ID         string
	Name       string
	Bytestream string
	MimeType   string
	Size       int64
}

type Message struct {
	ID        string
	Subject   string
	Timestamp string
	Files     []Attachment
}

type Code struct {
	CodeString string `json:"code"`
}
