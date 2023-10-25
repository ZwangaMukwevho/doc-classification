package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"doc-classification/pkg/resource"
	"doc-classification/pkg/service"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func main() {

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

	// initialise the gmail service
	ctx := context.Background()
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	// Setting up the user and the time stamp
	user := "me"
	currentTime := time.Now()
	yesterday := currentTime.AddDate(0, 0, -5)
	timestampTest := yesterday.Format("2006/01/02")
	query := fmt.Sprintf("in:inbox category:primary has:attachment after:%s", timestampTest)

	service.GetAttachmentArray(client, user, query, srv)
}
