package resource

import (
	"context"
	"doc-classification/pkg/common"
	"doc-classification/pkg/repository"
	"doc-classification/pkg/service"
	"errors"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/gmail/v1"
)

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config, tokenFile string) *http.Client {
	tok, err := loadOrRefreshToken(config, tokenFile)
	if err != nil {
		log.Fatal(err)
	}

	return config.Client(context.Background(), tok)
}

func GetClientFromDBToken(config *oauth2.Config, token *oauth2.Token, db repository.FirebaseRepository, userID string) (*http.Client, error) {

	if token.Valid() { // check if the token is expired
		common.Logger.Info("Error getting valid G-Auth token")
		return config.Client(context.Background(), token), nil
	}

	newToken, err := refreshToken(config, token, db, userID)
	if err != nil {
		common.Logger.Errorf("Error refreshing token %v", err)
		return nil, err
	}

	return config.Client(context.Background(), newToken), nil
}

func loadOrRefreshToken(config *oauth2.Config, tokenFile string) (*oauth2.Token, error) {
	tok, err := service.TokenFromFile(tokenFile)

	// If error is specifically about a expired token
	// if err != nil && err.Error() == "token is expired" {
	// 	return refreshAndSaveToken(config, tokenFile)
	// }

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

// func refreshAndSaveToken(config *oauth2.Config, tokenFile string) (*oauth2.Token, error) {

// 	// Reading the token
// 	f, err := os.Open(tokenFile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	tok := &oauth2.Token{}

// 	err = json.NewDecoder(f).Decode(tok)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Refresh the token
// 	newToken, err := refreshToken(config, tok)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Persist the new token
// 	service.SaveToken(tokenFile, newToken)
// 	return newToken, nil
// }

func refreshToken(config *oauth2.Config, token *oauth2.Token, db repository.FirebaseRepository, userID string) (*oauth2.Token, error) {
	ctx := context.Background()

	// Refresh the token
	tokenSource := config.TokenSource(ctx, token)

	if tokenSource == nil {
		return nil, errors.New("failed to get token source")
	}

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	if newToken == nil {
		return nil, errors.New("failed to get new token")
	}

	err = updateTokenInDB(config, db, *newToken, userID)
	if err != nil {
		return nil, errors.New("failed to update token in DB")
	}

	return newToken, nil
}

func updateTokenInDB(config *oauth2.Config, db repository.FirebaseRepository, token oauth2.Token, userID string) error {

	if config.Scopes[0] == gmail.GmailReadonlyScope {
		return db.UpdateGmailToken(userID, token)
	}

	if config.Scopes[0] == drive.DriveScope {
		return db.UpdateGdriveToken(userID, token)
	}

	// Update the token in the DB
	return nil
}
