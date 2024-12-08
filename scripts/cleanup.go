package scripts

import (
	"context"
	"doc-classification/pkg/repository"
	"doc-classification/pkg/resource"
	"doc-classification/pkg/service"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func CleanUpCron() {
	ctx := context.Background()
	b, err := os.ReadFile("google_client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	firebaseDB, err := repository.InitDB("https://react-getting-started-78f85-default-rtdb.firebaseio.com", "firebase_service.json")
	if err != nil {
		log.Fatal(err)
	}

	firebaseRepository := repository.NewFirebaseRestClient(firebaseDB)

	users, err := firebaseRepository.GetUserDataList()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("users %v: \n", users)

	for _, dbUserData := range *users {
		// Google drive setup
		driveConfig, err := google.ConfigFromJSON(b, drive.DriveFileScope)
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

		// Get current time and calculate the time for 2 days ago
		twoDaysAgo := time.Now().AddDate(0, 0, -2).Format(time.RFC3339)
		for categoryID, _ := range dbUserData.Categories {
			// List the files in the current category (assuming categoryID is a folder ID)
			r, err := localDriveService.Service.Files.List().
				Q("'" + categoryID + "' in parents and modifiedTime > '" + twoDaysAgo + "'").
				Fields("files(id, name, modifiedTime)").Do()

			if err != nil {
				log.Fatalf("Unable to retrieve files: %v", err)
			}

			for _, file := range r.Files {
				// Delete the file
				err := localDriveService.Service.Files.Delete(file.Id).Do()
				if err != nil {
					log.Printf("Unable to delete file: %v", err)
				} else {
					log.Printf("Deleted file: %s (%s)", file.Name, file.Id)
				}
			}
		}
	}

}
