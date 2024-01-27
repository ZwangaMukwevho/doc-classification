package service

import (
	"context"
	"doc-classification/pkg/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
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

// Retrieves a token from a local file.
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}

	if !tok.Valid() {  // check if the token is expired
		return nil, errors.New("token is expired")
	}
	
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		return nil, err
	}
	
	return tok, err
}

// Saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
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
