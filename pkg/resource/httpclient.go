package commmon

import (
	"context"
	"doc-classification/pkg/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config, tokenFile string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first time

	tokFile := tokenFile //TO-DO: have a unique one for each user
	tok, err := service.TokenFromFile(tokFile)
	fmt.Printf("Token: %+v\n", tok)
	fmt.Printf("Config: %+v\n", config)

	if err != nil {

		if err.Error() == "token is expired" { // if the token is expired
			// fetch a new token using the refresh token
			ctx := context.Background()

			// Reading the token
			f, err := os.Open(tokenFile)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			tok := &oauth2.Token{}

			err = json.NewDecoder(f).Decode(tok)
			if err != nil {
				log.Fatal(err)
			}

			// Refreshing
			tokenSource := config.TokenSource(ctx, tok)
			if tokenSource == nil {
				log.Fatal("Failed to get token source")
			}
			newToken, err := tokenSource.Token()
			if err != nil {
				log.Fatal(err)
			} else if newToken == nil {
				log.Fatal("Failed to get new token")
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
	fmt.Printf("Token at the end: %+v\n", tok)
	return config.Client(context.Background(), tok)
}
