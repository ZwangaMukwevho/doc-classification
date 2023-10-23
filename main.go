package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"doc-classification/pkg/resource"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	b, err := os.ReadFile("client_secret_973692223612-28ae9a7njdsfh7gv89l0fih5q36jt52m.apps.googleusercontent.com.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := resource.GetClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	// r, err := srv.Users.Labels.List(user).Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve labels: %v", err)
	// }
	// if len(r.Labels) == 0 {
	// 	fmt.Println("No labels found.")
	// 	return
	// }
	// fmt.Println("Labels:")
	// for _, l := range r.Labels {
	// 	fmt.Printf("- %s\n", l.Name)
	// }

	messages, err := srv.Users.Messages.List(user).Q("category:primary").MaxResults(8).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve unread messages: %v", err)
	}

	fmt.Println("Subject headings of the 5 most recent emails:")
	for _, message := range messages.Messages {
		msg, err := srv.Users.Messages.Get(user, message.Id).Do()
		if err != nil {
			log.Printf("Unable to retrieve message details: %v", err)
			continue
		}
		headers := msg.Payload.Headers
		for _, header := range headers {
			if header.Name == "Subject" {
				fmt.Printf("- %s\n", header)
				break
			}
		}

		// Check for attachments
		parts := msg.Payload.Parts
		if parts != nil {
			for _, part := range parts {
				if part.Filename != "" {
					fmt.Printf("  - Attachment Filename: %s\n", part.Filename)
				}
			}
		}

	}

}
