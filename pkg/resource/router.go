package resource

import (
	"context"
	"doc-classification/pkg/service"
	"net/http"

	"golang.org/x/oauth2"
)

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config, tokenFile string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.

	tokFile := tokenFile // to-do: have a unique one for each user
	tok, err := service.TokenFromFile(tokFile)
	if err != nil {
		tok = service.GetTokenFromWeb(config)
		service.SaveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}
