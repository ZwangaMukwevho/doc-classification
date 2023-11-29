package main

import (
	"context"
	"doc-classification/pkg/resource"
	"doc-classification/pkg/service"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

/*
	Purposefully comment out the below code
	Only uncomment when you want to get the directory ID's from your google drive
	Only needed this while setting up the directories.json file
*/
// dirId, err := localDriveService.GetDriveDirectories()

func main() {
	ctx := context.Background()
	oAuthFileName := "client_secret_973692223612-28ae9a7njdsfh7gv89l0fih5q36jt52m.apps.googleusercontent.com.json"

	b, err := os.ReadFile(oAuthFileName)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	//Gmail Setup
	gmailConfig, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	resource.GetClient(gmailConfig, "token_gmail.json")

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

	// Setting up the directories
	localDriveService := service.DriveServiceLocal{Service: driveSrv}

	/*
		Purposefully comment out the below code
		Only uncomment when you want to get the directory ID's from your google drive
		Only needed this while setting up the directories.json file
	*/
	dirId, err := localDriveService.GetDriveDirectories()
	fmt.Println(dirId)

	file, err := os.Create("output.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Create a JSON encoder
	encoder := json.NewEncoder(file)

	// Encode the list of directories to JSON and write to the file
	err = encoder.Encode(dirId)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println("JSON file created successfully.")

}
