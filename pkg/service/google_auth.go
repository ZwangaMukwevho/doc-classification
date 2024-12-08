package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"doc-classification/pkg/common"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuthMethods interface {
	GetTokenFromWeb(config *oauth2.Config) *oauth2.Token
	GetTokenUsingAPI(config *oauth2.Config, code string) (*oauth2.Token, error)
	GetAuthCodeURL(config *oauth2.Config) string
	TokenFromFile(file string) (*oauth2.Token, error)
}

func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	common.Logger.Infof("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		common.Logger.Errorf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		common.Logger.Errorf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func GetTokenUsingAPI(config *oauth2.Config, code string) (*oauth2.Token, error) {
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		common.Logger.Errorf("Unable to retrieve token from web: %v", err)
		return nil, err
	}

	return tok, nil
}

func GetAuthCodeURL(config *oauth2.Config) string {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return authURL
}

func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		common.Logger.Errorf("Token file not found: %v", err)
		return nil, errors.New("Token file not found")
	}

	defer f.Close()
	tok := &oauth2.Token{}

	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		common.Logger.Errorf("Error decoding token: %v", err)
		return nil, err
	}

	if !tok.Valid() { // check if the token is expired
		common.Logger.Warnf("Token from file %s has expired: ", file)
		return nil, errors.New("token is expired")
	}

	return tok, err
}

func SaveToken(path string, token *oauth2.Token) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		common.Logger.Errorf("Unable to cache oauth token: %v", err)
	}

	defer file.Close()
	json.NewEncoder(file).Encode(token)
}

func GetGoogleToken(code string, scope string) (*oauth2.Token, error) {
	config, err := GetOauthConfig(scope)
	if err != nil {
		common.Logger.Errorf("Error getting Oauth Config: %v", err)
		return nil, err
	}

	token, err := GetTokenUsingAPI(config, code)
	if err != nil {
		common.Logger.Errorf("Error retrieving gmail token: %v", err)
		return nil, err
	}

	return token, nil
}

func GetOauthConfig(scope string) (*oauth2.Config, error) {
	oAuthByteStream, err := common.GetJsonFileByteStream("google_client_secret.json")
	if err != nil {
		common.Logger.Errorf("Error reading client secret file: %v", err)
		return nil, err
	}

	config, err := google.ConfigFromJSON(*oAuthByteStream, scope)
	if err != nil {
		common.Logger.Errorf("Errpr getting config from client credentials files: %v", err)
		return nil, err
	}

	return config, nil
}
