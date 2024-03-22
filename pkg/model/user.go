package model

import "golang.org/x/oauth2"

type User struct {
	UserId     string     `json:"userid"`
	Categories []Category `json:"categories"`
	GmailCode  string     `json:"gmailCode"`
	GdriveCode string     `json:"gdriveCode"`
}

type FirebaseUser struct {
	UserId     string
	Categories map[string]Category
	GmailCode  *oauth2.Token
	GdriveCode *oauth2.Token
}

type Users struct {
	UserId     string              `json:"UserId"`
	Categories map[string]Category `json:"Categories"`
	GmailCode  *oauth2.Token       `json:"GmailCode"`
	GdriveCode *oauth2.Token       `json:"GdriveCode"`
}

type Category struct {
	Category    string `json:"category"`
	Description string `json:"description"`
}
