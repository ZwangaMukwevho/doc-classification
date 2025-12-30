package cron

import (
	"context"
	"doc-classification/pkg/common"
	"doc-classification/pkg/gateway"
	"doc-classification/pkg/repository"
	"doc-classification/pkg/resource"
	"doc-classification/pkg/service"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func ClassificationCron() {
	common.Logger.Info("Cron job is running at: ", time.Now())

	// Add your cron job logic here
	if err := godotenv.Load(); err != nil {
		common.Logger.Error("Error loading .env file")
		return
	}

	// Read API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		common.Logger.Error("API key is missing. Set the OPENAI_API_KEY environment variable.")
		return
	}

	ctx := context.Background()
	b, err := os.ReadFile("google_client_secret.json")
	if err != nil {
		common.Logger.Fatalf("Unable to read client secret file: %v", err)
	}

	firebaseDB, err := repository.InitDB("https://react-getting-started-78f85-default-rtdb.firebaseio.com", "firebase_service.json")
	if err != nil {
		common.Logger.Fatalf("Unable to read client secret file from firebase: %v", err)
	}

	// Getting user data with tokens from DB
	firebaseRepository := repository.NewFirebaseRestClient(firebaseDB)
	users, err := firebaseRepository.GetUserDataList()
	if err != nil {
		common.Logger.Errorf("Unable to read users from firebase: %v", err)
	}

	for _, dbUserData := range *users {
		common.Logger.Infof("Processing user %v", dbUserData.UserId)

		gmailConfig, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
		if err != nil {
			common.Logger.Errorf("Unable to parse client secret file to config: %v", err)
			continue
		}

		driveConfig, err := google.ConfigFromJSON(b, drive.DriveFileScope)
		if err != nil {
			common.Logger.Errorf("Unable to parse client secret file to config: %v", err)
			continue
		}

		gmailClient, err := resource.GetClientFromDBToken(gmailConfig, dbUserData.GmailCode, firebaseRepository, dbUserData.UserId)
		common.Logger.Infof("Gmail client %v", gmailClient)
		if err != nil {
			common.Logger.Errorf("Unable to get gmail client: %v", err)
			continue
		}

		gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(gmailClient))
		if err != nil {
			common.Logger.Errorf("Unable to retrieve Gmail client: %v", err)
			continue

		}
		localGmailService := service.GmailServiceLocal{Service: gmailService}

		// dereferencedMessageArr := *messagesArray
		driveClient, err := resource.GetClientFromDBToken(driveConfig, dbUserData.GdriveCode, firebaseRepository, dbUserData.UserId)
		if err != nil {
			common.Logger.Errorf("Unable to get gmail client: %v", err)
			continue
		}

		driveService, err := drive.NewService(ctx, option.WithHTTPClient(driveClient))
		if err != nil {
			common.Logger.Errorf("Unable to retrieve Drive client: %v", err)
			continue
		}

		localDriveService := service.DriveServiceLocal{Service: driveService}

		// Query the attachments
		user := "me"
		queryDateRange := time.Now().AddDate(0, 0, -1).Format("2006/01/02")
		common.Logger.Infof("query date currentyl is %v", queryDateRange)
		query := fmt.Sprintf("in:inbox category:primary has:attachment after:%s -from:no-reply@sixty60.co.za", queryDateRange)

		messagesArray, err := localGmailService.GetAttachmentArray(user, query)
		if err != nil {
			common.Logger.Errorf("error getting the attachments: %v", err)
			continue
		}

		/*
			Purposefully comment out the below code
			Only uncomment when you want to get the directory ID's from your google drive
			Only needed this while setting up the directories.json file
		*/
		//dirId, err := localDriveService.GetDriveDirectories()

		// Get the file directories and ID's

		openAIContentString := service.CreateContentString(dbUserData.Categories)
		for _, message := range *messagesArray {
			// Create the classification prompt
			subject := message.Subject
			for _, attachment := range message.Files {
				classificationQuestion := service.CreateClassifyEmailPrompt(subject, attachment)
				classificationPrompt := service.CreateSubsequentPrompt(classificationQuestion)

				// Send the classification request
				classificationResponse, err := gateway.SendCompletionRequest(openAIContentString, classificationPrompt, apiKey)
				if err != nil {
					common.Logger.Errorf("Error sending classification request to openai with : %v", err)
					continue
				}

				if classificationResponse == nil {
					continue
				}

				oneWordResponse, err := service.ExtractOpenAIContent(*classificationResponse)
				if err != nil {
					common.Logger.Errorf("Error extracting response from openai : %v", err)
					continue
				}

				driveDirID, err := common.FindDirectoryByID(dbUserData.Categories, *oneWordResponse)
				if err != nil {
					common.Logger.Errorf("Error getting corresponding google drive id locally : %v", err)
					continue
				}

				fileExists := localDriveService.FileExists(attachment.Name, *driveDirID)
				if fileExists != nil {
					common.Logger.Errorf("Skipped uploading file due to the following error: %v", fileExists)
					continue
				}

				// Finally upload file
				driveUploadErr := localDriveService.UploadFile(attachment, *driveDirID)
				if driveUploadErr != nil {
					common.Logger.Errorf("error uploading file to gdrive %v \n", err)
					continue
				}
			}
		}
	}

	common.Logger.Info("Cron job is finished running at: ", time.Now())
}
