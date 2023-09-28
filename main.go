package main

import (
    "context"
    "fmt"
    "log"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    gmail "google.golang.org/api/gmail/v1"
)

func main() {
    // Load your OAuth 2.0 credentials from the JSON file
    credentialsFile := "path/to/your/credentials.json"
    ctx := context.Background()
    config, err := google.ConfigFromJSONFile(credentialsFile, gmail.GmailReadonlyScope)
    if err != nil {
        log.Fatalf("Unable to parse client secret file: %v", err)
    }

    // Get an OAuth2 token
    client := getClient(ctx, config)

    // Create a Gmail API client
    srv, err := gmail.New(client)
    if err != nil {
        log.Fatalf("Unable to retrieve Gmail client: %v", err)
    }

    // You can now use 'srv' to interact with the Gmail API.
    // For example, you can list the user's labels:
    labels, err := srv.Users.Labels.List("me").Do()
    if err != nil {
        log.Fatalf("Unable to retrieve labels: %v", err)
    }
    fmt.Println("Labels:")
    for _, label := range labels.Labels {
        fmt.Printf("- %s\n", label.Name)
    }
}

func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
    tokenFile := "path/to/your/token.json" // You may want to store the token for later use
    tok, err := tokenFromFile(tokenFile)
    if err != nil {
        tok = getTokenFromWeb(config)
        saveToken(tokenFile, tok)
    }
    return config.Client(ctx, tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
    // Load a token from a file (if it exists)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
    // Obtain an OAuth 2.0 token interactively (for initial setup)
}

func saveToken(file string, token *oauth2.Token) {
    // Save a token to a file (for future use)
}
