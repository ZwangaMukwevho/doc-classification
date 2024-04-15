package main

import (
	"context"
	"doc-classification/pkg/common"
	"doc-classification/pkg/gateway"
	"doc-classification/pkg/repository"
	"doc-classification/pkg/resource"
	"doc-classification/pkg/service"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func main() {

	// Load environment variables
	// Load environment variables from .env file
	go setupCron()

	firebaseDB, err := repository.InitDB("https://react-getting-started-78f85-default-rtdb.firebaseio.com", "firebase_service.json")
	if err != nil {
		log.Fatal(err)
	}

	basePath := "localhost:8080"

	firebaseRepository := repository.NewFirebaseRestClient(firebaseDB)

	router := resource.NewRouter(
		resource.Handler{
			FirebaseClient:      firebaseDB,
			FirebaseRespository: firebaseRepository,
		},
	)

	router.Run(basePath)
}

func cronJob() {
	fmt.Println("Cron job is running at:", time.Now())
	// Add your cron job logic here
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
		return
	}

	// Read API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("API key is missing. Set the OPENAI_API_KEY environment variable.")
		return
	}

	ctx := context.Background()
	b, err := os.ReadFile("client_secret_973692223612-28ae9a7njdsfh7gv89l0fih5q36jt52m.apps.googleusercontent.com.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	firebaseDB, err := repository.InitDB("https://react-getting-started-78f85-default-rtdb.firebaseio.com", "firebase_service.json")
	if err != nil {
		log.Fatal(err)
	}

	// Getting user data with tokens from DB
	firebaseRepository := repository.NewFirebaseRestClient(firebaseDB)
	users, err := firebaseRepository.GetUserDataList()
	if err != nil {
		log.Fatal(err)
	}

	for _, dbUserData := range *users {
		gmailConfig, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}

		gmailClient, err := resource.GetClientFromDBToken(gmailConfig, dbUserData.GmailCode, firebaseRepository, dbUserData.UserId)
		if err != nil {
			log.Fatalf("Unable to get gmail client: %v", err)
		}

		// initialise the gmail service
		srv, err := gmail.NewService(ctx, option.WithHTTPClient(gmailClient))
		if err != nil {
			log.Fatalf("Unable to retrieve Gmail client: %v", err)
		}

		// Setting up the user and the time stamp
		user := "me"
		currentTime := time.Now()
		yesterday := currentTime.AddDate(0, 0, -1)
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

		driveClient, err := resource.GetClientFromDBToken(driveConfig, dbUserData.GdriveCode, firebaseRepository, dbUserData.UserId)
		if err != nil {
			log.Fatalf("Unable to get gmail client: %v", err)
		}

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

		openAIContentString := service.CreateContentString(dbUserData.Categories)
		for _, message := range dereferencedMessageArr {
			// Create the classification prompt
			subject := message.Subject
			for _, attachment := range message.Files {
				classificationQuestion := service.CreateClassifyEmailPrompt(subject, attachment)
				classificationPrompt := service.CreateSubsequentPrompt(classificationQuestion)

				// Send the classification request
				classificationResponse, err := gateway.SendCompletionRequest(openAIContentString, classificationPrompt, apiKey)
				if err != nil {
					log.Fatalf("Error sending classification request to openai with : %v", err)
				}

				fmt.Printf("email subject name: %s , email attachment name: %s \n", message.Subject, attachment.Name)
				if classificationResponse != nil {
					oneWordResponse, err1 := service.ExtractOpenAIContent(*classificationResponse)
					if err1 != nil {
						log.Print("Error extracting response from")
					}

					fmt.Printf("Formatted string: %s \n", *oneWordResponse)
					driveDirID, err := common.FindDirectoryByID(dbUserData.Categories, *oneWordResponse)
					if err != nil {
						log.Fatalf("Error getting corresponding google drive id locally : %v", err)
					}

					// Finally upload file
					driveUploadErr := localDriveService.UploadFile(attachment, *driveDirID)
					if driveUploadErr != nil {
						log.Printf("error uploading file %v \n", err)
					}
				}
			}
		}
	}
}

func setupCron() {
	c := cron.New()

	// Schedule the job to run every minute
	// */3 * * * * fixing
	// 0 0 * * * normal
	_, err := c.AddFunc("0 0 * * *", cronJob)
	if err != nil {
		log.Println("Error scheduling cron job:", err)
		return
	}

	c.Start()

	// Run the cron scheduler in the background
	select {}
}
