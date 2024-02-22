package gateway

import (
	"doc-classification/pkg/resource"
	"log"
	"net/http"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func GetGmailClient(ouathFile *[]byte) *http.Client {
	// Read in the oAuthFile
	gmailConfig, err := google.ConfigFromJSON(*ouathFile, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	gmailClient := resource.GetClient(gmailConfig, "token_gmail.json")
	return gmailClient
}
