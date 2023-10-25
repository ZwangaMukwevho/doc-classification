package service

import (
	"context"
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
	err = json.NewDecoder(f).Decode(tok)
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

func GetAttachmentArray(client *http.Client, user string, query string, service *gmail.Service) {

	messages, err := service.Users.Messages.List(user).Q(query).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve unread messages: %v", err)
	}

	// Loop through all the messages array
	for _, message := range messages.Messages {
		msg, err := service.Users.Messages.Get(user, message.Id).Do()
		if err != nil {
			log.Printf("Unable to retrieve message details: %v", err)
			continue
		}
		headers := msg.Payload.Headers
		msg_id := message.Id
		timestamp := ""
		fmt.Printf("- %s\n:", msg_id)
		for _, header := range headers {

			if header.Name == "Date" {
				timestamp = header.Value
			}

			if header.Name == "Subject" {
				fmt.Printf("Timestamp: %s\n", timestamp)
				fmt.Printf("- Subject: %s\n", header)
				break
			}
		}

		// Check for attachments
		parts := msg.Payload.Parts
		if parts != nil {
			for _, part := range parts {
				if part.Filename != "" {
					attachmentID := part.Body.AttachmentId
					if attachmentID != "" {
						attachmentLink := fmt.Sprintf("https://mail.google.com/mail/u/0?ik=YOUR_USER_ID&attid=%s", attachmentID)
						fmt.Printf("  - Attachment Link: %s\n", attachmentLink)
					}
					fmt.Printf("  - Attachment Filename: %s\n", part.Filename)
				}
			}
		}
	}
}
