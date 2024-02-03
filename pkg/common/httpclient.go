package common

import (
	"context"
	"doc-classification/pkg/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config, tokenFile string) *http.Client {
	tok, err := loadOrRefreshToken(config, tokenFile)
	if err != nil {
		log.Fatal(err)
	}

	return config.Client(context.Background(), tok)
}

func loadOrRefreshToken(config *oauth2.Config, tokenFile string) (*oauth2.Token, error) {
	tok, err := service.TokenFromFile(tokenFile)

	// If error is specifically about a expired token
	if err != nil && err.Error() == "token is expired" {
		return refreshAndSaveToken(config, tokenFile)
	}

	if err != nil && err.Error() == "Token file not found" {
		tok = service.GetTokenFromWeb(config)
		service.SaveToken(tokenFile, tok)
		return tok, nil
	}

	// General load
	if err != nil {
		return nil, err
	}
	return tok, nil
}

func refreshAndSaveToken(config *oauth2.Config, tokenFile string) (*oauth2.Token, error) {
	ctx := context.Background()

	// Reading the token
	f, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}

	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		log.Fatal(err)
	}

	// Refresh the token
	tokenSource := config.TokenSource(ctx, tok)

	if tokenSource == nil {
		return nil, errors.New("failed to get token source")
	}

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	} else if newToken == nil {
		return nil, errors.New("failed to get new token")
	}

	// Persist the new token
	service.SaveToken(tokenFile, newToken)
	return newToken, nil
}
