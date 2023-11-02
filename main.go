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
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	b, err := os.ReadFile("client_secret_973692223612-28ae9a7njdsfh7gv89l0fih5q36jt52m.apps.googleusercontent.com.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	//If modifying these scopes, delete your previously saved token.json.
	//Gmail Setup
	gmailConfig, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	gmailClient := resource.GetClient(gmailConfig, "token_gmail.json")

	// initialise the gmail service
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(gmailClient))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	// Setting up the user and the time stamp
	user := "me"
	currentTime := time.Now()
	yesterday := currentTime.AddDate(0, 0, -10)
	timestampTest := yesterday.Format("2006/01/02")
	query := fmt.Sprintf("in:inbox category:primary has:attachment after:%s", timestampTest)

	messagesArray, err := service.GetAttachmentArray(gmailClient, user, query, srv)
	if err != nil {
		log.Print("error getting the attachments")
	}

	// Google drive setup
	driveConfig, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	driveClient := resource.GetClient(driveConfig, "token_g_drive.json")
	driveSrv, err := drive.NewService(ctx, option.WithHTTPClient(driveClient))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	r, err1 := driveSrv.Files.List().PageSize(10).Fields("files(id, name)").Do()
	if err1 != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	// fmt.Println("Files:")
	// fmt.Println(r)
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}

	attachment1 := *messagesArray
	//mimeType := common.GetMimeType(attachment1[0].File.Name)
	fmt.Println("mime type")
	fmt.Println(attachment1[0].File.MimeType)
	// targetDirectoryID := "test"
	// file := &drive.File{
	// 	Name:     attachment1[0].File.Name,
	// 	MimeType: attachment1[0].File.MimeType,
	// 	Parents:  []string{"test_dir"},
	// }

	// // Decode the base64url data
	// data, err := base64.RawURLEncoding.DecodeString(attachment1[0].File.Bytestream)
	// if err != nil {
	// 	log.Fatalf("Error decoding base64url data: %v", err)
	// }

	// _, err = driveSrv.Files.Create(file).Media(bytes.NewReader((data))).Do()
	// if err != nil {
	// 	log.Fatalf("Unable to create file: %v", err)
	// }
	// fmt.Printf("File '%s' uploaded to the specified directory in Google Drive.\n", attachment1[0].File.Name)

}
