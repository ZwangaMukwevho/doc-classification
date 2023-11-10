package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// ctx := context.Background()
	// b, err := os.ReadFile("client_secret_973692223612-28ae9a7njdsfh7gv89l0fih5q36jt52m.apps.googleusercontent.com.json")
	// if err != nil {
	// 	log.Fatalf("Unable to read client secret file: %v", err)
	// }

	// //If modifying these scopes, delete your previously saved token.json.
	// //Gmail Setup
	// gmailConfig, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	// if err != nil {
	// 	log.Fatalf("Unable to parse client secret file to config: %v", err)
	// }
	// gmailClient := resource.GetClient(gmailConfig, "token_gmail.json")

	// // initialise the gmail service
	// srv, err := gmail.NewService(ctx, option.WithHTTPClient(gmailClient))
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve Gmail client: %v", err)
	// }

	// // Setting up the user and the time stamp
	// user := "me"
	// currentTime := time.Now()
	// yesterday := currentTime.AddDate(0, 0, -10)
	// timestampTest := yesterday.Format("2006/01/02")
	// query := fmt.Sprintf("in:inbox category:primary has:attachment after:%s", timestampTest)

	// messagesArray, err := service.GetAttachmentArray(gmailClient, user, query, srv)
	// if err != nil {
	// 	log.Print("error getting the attachments")
	// }

	// // Google drive setup
	// driveConfig, err := google.ConfigFromJSON(b, drive.DriveScope)
	// if err != nil {
	// 	log.Fatalf("Unable to parse client secret file to config: %v", err)
	// }
	// driveClient := resource.GetClient(driveConfig, "token_g_drive.json")
	// driveSrv, err := drive.NewService(ctx, option.WithHTTPClient(driveClient))
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve Drive client: %v", err)
	// }

	// localDriveService := service.DriveServiceLocal{Service: driveSrv}
	// attachment1 := *messagesArray
	// err = localDriveService.UploadFile(attachment1[0])

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
	apiUrl := "https://api.openai.com/v1/chat/completions"

	// Request payload
	requestBody := `{
        "model": "gpt-3.5-turbo",
        "messages": [
            {"role": "system", "content": "You are a document classification assistant. I have categories I want to classify my email documents which are: 1. Education: This category can include documents related to educational pursuits, such as school transcripts, certificates, course materials, and research papers. 2. Finance: Finance-related documents can cover a wide range of items, including bank statements, tax records, invoices, quotes, receipts, and investment reports. 5. Work: Work-related documents can involve project plans, reports, emails, resumes, and other materials directly related to your professional life. 6. Home: Home category files may include property documents, utility bills, home maintenance records, and household inventory. 7. Personal: This category can cover a wide range of personal documents, from family photos to personal notes, travel itineraries, and more. Give your reply as one word answer from the given categories"},
            {"role": "user", "content": "How can you classify attachment from an email with subject 'Quotation:21026480' and attachment name: 'Vehicle Booking Acceptance form ref::2102648'"}
        ]
    }`

	// Create a request
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBufferString(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the API key in the request headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and display the response
	responseBody := new(bytes.Buffer)
	responseBody.ReadFrom(resp.Body)
	fmt.Println("Response:", responseBody.String())
}
