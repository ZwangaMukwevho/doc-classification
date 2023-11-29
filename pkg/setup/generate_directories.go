package setup

import (
	"context"
	"doc-classification/pkg/resource"
	"doc-classification/pkg/service"
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

	relativeOauthFileName := "../" + oAuthFileName
	b, err := os.ReadFile(relativeOauthFileName)
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
	fmt.Println(localDriveService)
	/*
		Purposefully comment out the below code
		Only uncomment when you want to get the directory ID's from your google drive
		Only needed this while setting up the directories.json file
	*/
	//dirId, err := localDriveService.GetDriveDirectories()

}
