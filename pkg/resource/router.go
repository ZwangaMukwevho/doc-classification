package resource

import (
	"context"
	"doc-classification/pkg/service"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config, tokenFile string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first time

	tokFile := tokenFile //TO-DO: have a unique one for each user
	tok, err := service.TokenFromFile(tokFile)

	if err != nil {

		if err.Error() == "token is expired" { // if the token is expired
			// fetch a new token using the refresh token
			ctx := context.Background()
			tokenSource := config.TokenSource(ctx, tok)
			newToken, err := tokenSource.Token()
			if err != nil {
				log.Fatal(err)
			}
			// persist the new token
			service.SaveToken(tokFile, newToken)
			// use the new token
			tok = newToken // set the expired token to the new token
		} else {
			tok = service.GetTokenFromWeb(config)
			service.SaveToken(tokFile, tok)
		}
	}
	return config.Client(context.Background(), tok)
}
