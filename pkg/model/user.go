package model

import "golang.org/x/oauth2"

type User struct {
	UserId     string   `json:"userid"`
	Categories []string `json:"categories"`
	GmailCode  string   `json:"gmailCode"`
	GdriveCode string   `json:"gdriveCode"`
}

type FirebaseUser struct {
	UserId     string
	Categories map[string]string
	GmailCode  *oauth2.Token
	GdriveCode *oauth2.Token
}

type Users struct {
	UserId     string        `json:"UserId"`
	Categories []string      `json:"Categories"`
	GmailCode  *oauth2.Token `json:"GmailCode"`
	GdriveCode *oauth2.Token `json:"GdriveCode"`
}
