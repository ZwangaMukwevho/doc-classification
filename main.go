package main

import (
	"context"
	"doc-classification/pkg/common"
	"doc-classification/pkg/gateway"
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

	// Load environment variables
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
	dereferencedMessageArr := *messagesArray

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

	/*
		Purposefully comment out the below code
		Only uncomment when you want to get the directory ID's from your google drive
		Only needed this while setting up the directories.json file
	*/
	//dirId, err := localDriveService.GetDriveDirectories()

	// Get the file directories and ID's
	directories, fileErr := common.ReadJsonFile("pkg/common/directories.json")
	if fileErr != nil {
		log.Panicf("Error reading file directories %v: ", *fileErr)
	}

	for _, message := range dereferencedMessageArr {
		// Create the classification prompt
		subject := message.Subject
		for _, attachment := range message.Files {
			classificationQuestion := service.CreateClassifyEmailPrompt(subject, attachment)
			classificationPrompt := service.CreateSubsequentPrompt(classificationQuestion)

			// Send the classification request
			classificationResponse, err := gateway.SendCompletionRequest(classificationPrompt, apiKey)
			if err != nil {
				log.Fatalf("Error sending classification request to openai with : %v", err)
			}

			fmt.Println("file info")
			fmt.Printf("email subject name: %s , email attachment name: %s \n", message.Subject, attachment.Name)

			if classificationResponse != nil {
				oneWordResponse, err1 := service.ExtractOpenAIContent(*classificationResponse)
				if err1 != nil {
					log.Print("Error extracting response from")
				}

				fmt.Printf("Formatted string: %s \n", *oneWordResponse)
				driveDirID, err := common.FindDirectoryByID(*directories, *oneWordResponse)
				if err != nil {
					log.Fatalf("Error getting corresponding google drive id locally : %v", err)
				}
				fmt.Printf("Corresponding file id: %s \n", *driveDirID)

				// Finally upload file
				driveUploadErr := localDriveService.UploadFile(attachment, *driveDirID)
				if driveUploadErr != nil {
					fmt.Printf("error uploading file %v \n", err)
				}

			}
		}
	}

}
