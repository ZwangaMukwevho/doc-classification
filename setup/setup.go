package main

import (
	"context"
	"doc-classification/pkg/resource"
	"doc-classification/pkg/service"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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
	// Add your cron job logic here
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Read API key from environment variable
	gDriveTokenFile := os.Getenv("G_DRIVE_TOKEN_FILE")
	gmailTokenFile := os.Getenv("GMAIL_TOKEN_FILE")
	oAuthFileName := os.Getenv("GOOGLE_AUTH_FILE")

	ctx := context.Background()
	oAuthFileRelative := "../" + oAuthFileName

	b, err := os.ReadFile(oAuthFileRelative)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	//Gmail Setup
	gmailConfig, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	gmailTokenRelativePath := "../" + gmailTokenFile
	resource.GetClient(gmailConfig, gmailTokenRelativePath)

	// Google drive setup
	driveConfig, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	gdriveTokenFileRelativePath := "../" + gDriveTokenFile
	driveClient := resource.GetClient(driveConfig, gdriveTokenFileRelativePath)
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

	file, err := os.Create("../pkg/common/directories.json")
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
