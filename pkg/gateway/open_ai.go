package gateway

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func SendCompletionRequest(contentString string, userRequestBody string, apiKey string) (*string, error) {
	apiUrl := "https://api.openai.com/v1/chat/completions"

	requestBody := fmt.Sprintf(`{
	    "model": "gpt-3.5-turbo",
	    "messages": [
	        {"role": "system", "content": "%s"},
	        %s
	    ]
	}`, contentString, userRequestBody)

	// Create a request
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBufferString(requestBody))
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	// Set the API key in the request headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()
	// Read and display the response
	responseBody := new(bytes.Buffer)
	responseBody.ReadFrom(resp.Body)

	responseString := responseBody.String()
	return &responseString, nil
}
