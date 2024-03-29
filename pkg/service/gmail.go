package service

import (
	"context"
	"doc-classification/pkg/common"
	"doc-classification/pkg/model"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Request a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func GetTokenUsingAPI(config *oauth2.Config, code string) (*oauth2.Token, error) {
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Printf("Unable to retrieve token from web: %v", err)
		return nil, err
	}
	return tok, nil
}

func GetAuthCodeURL(config *oauth2.Config) string {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL
}

// Retrieves a token from a local file.
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.New("Token file not found")
	}
	defer f.Close()
	tok := &oauth2.Token{}

	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		return nil, err
	}

	if !tok.Valid() { // check if the token is expired
		return nil, errors.New("token is expired")
	}

	return tok, err
}

// Saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetAttachmentArray(client *http.Client, user string, query string, service *gmail.Service) (*[]model.Message, error) {
	messages, err := service.Users.Messages.List(user).Q(query).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve unread messages: %v", err)
		return nil, err
	}

	// Initialise variables
	var messagesArray []model.Message
	var messageStruct model.Message

	// Loop through all the messages array
	for _, message := range messages.Messages {

		// Get messages
		msg, err := service.Users.Messages.Get(user, message.Id).Do()
		if err != nil {
			log.Printf("Unable to retrieve message details: %v", err)
			continue
		}

		headers := msg.Payload.Headers
		timestamp := ""
		for _, header := range headers {

			if header.Name == "Date" {
				timestamp = header.Value
				messageStruct.Timestamp = timestamp
			}

			if header.Name == "Subject" {
				messageStruct.Subject = header.Value
				break
			}
		}

		// Check for attachments
		parts := msg.Payload.Parts
		var attachments []model.Attachment

		if parts != nil {
			for _, part := range parts {
				if part.Filename != "" {
					attachmentID := part.Body.AttachmentId
					// attachment.ID = attachmentID
					// attachment.Name = part.Filename
					// Get the attachment Bytestream
					if attachmentID != "" {
						attachmentData, err := service.Users.Messages.Attachments.Get(user, message.Id, attachmentID).Do()
						if err != nil {
							log.Printf("Unable to retrieve attachment content: %v", err)
							continue
						}

						attachment := model.Attachment{
							ID:         attachmentID,
							Name:       part.Filename,
							MimeType:   part.MimeType,
							Bytestream: attachmentData.Data,
							Size:       attachmentData.Size,
						}

						attachments = append(attachments, attachment)
						// attachment.Bytestream = attachmentData.Data
						// attachment.Size = attachmentData.Size
					}
				}
				// attachment.MimeType = part.MimeType
			}
		}

		messageStruct.Files = attachments
		messagesArray = append(messagesArray, messageStruct)
	}
	return &messagesArray, nil
}

func GetGmailOauthConfig(scope string) (*oauth2.Config, error) {
	oAuthByteStream, err := common.GetJsonFileByteStream("client_secret_973692223612-28ae9a7njdsfh7gv89l0fih5q36jt52m.apps.googleusercontent.com.json")
	if err != nil {
		return nil, err
	}

	gmailConfig, err := google.ConfigFromJSON(*oAuthByteStream, scope)
	if err != nil {
		return nil, err
	}

	return gmailConfig, nil
}

func GetGmailToken(code string) (*oauth2.Token, error) {
	config, err := GetGmailOauthConfig(gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}

	gmailToken, err := GetTokenUsingAPI(config, code)
	if err != nil {
		return nil, err
	}

	return gmailToken, nil
}
