package main

import (
	"context"
	"doc-classification/pkg/resource"
	"doc-classification/pkg/service"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
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
	query := fmt.Sprintf("in:inbox category:primary has:attachment after:%s -from:no-reply@sixty60.co.za", timestampTest)

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

	localDriveService := service.DriveServiceLocal{Service: driveSrv}
	attachment1 := *messagesArray

	for i, a := range attachment1 {

		// if i == 1 {
		fmt.Printf("file name %d: ", i)
		fmt.Println(a.File.Name)
		fmt.Println(a.File.Size)

		err = localDriveService.UploadFile(a)
		if err != nil {
			fmt.Printf("error uploading file %v \n", err)
		}

		// Specify the file path
		// filePath := "output.txt"
		// // Write the string to the file
		// _ = ioutil.WriteFile(filePath, []byte(a.File.Bytestream), 0644)
		// fmt.Printf("mime type %v \n", a.File.MimeType)
		// }

	}
	//err = localDriveService.UploadFile(attachment1[0])

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Read API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("API key is missing. Set the OPENAI_API_KEY environment variable.")
		return
	}

	// response, err := gateway.SendCompletionRequest("test", apiKey)
	// if err != nil {
	// 	fmt.Printf("error %v: ", err)
	// }
	// fmt.Println("Response:", *response)
}
