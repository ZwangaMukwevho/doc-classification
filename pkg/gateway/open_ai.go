package gateway

import (
	"bytes"
	"fmt"
	"net/http"
)

func SendCompletionRequest(userRequestBody string, apiKey string) (*string, error) {
	apiUrl := "https://api.openai.com/v1/chat/completions"

	requestBody := fmt.Sprintf(`{
	    "model": "gpt-3.5-turbo",
	    "messages": [
	        {"role": "system", "content": "You are a document classification assistant. I have categories I want to classify my email documents which are: 1. Education: This category can include documents related to educational pursuits, such as school transcripts, certificates, course materials, and research papers. 2. Finance: Finance-related documents can cover a wide range of items, including bank statements, tax records, invoices, quotes, receipts, and investment reports. 5. Work: Work-related documents can involve project plans, reports, emails, resumes, and other materials directly related to your professional life. 6. Home: Home category files may include property documents, utility bills, home maintenance records, and household inventory. 7. Personal: This category can cover a wide range of personal documents, from family photos to personal notes, travel itineraries, and more. Give your reply as one word answer from the given categories"},
	        %s
	    ]
	}`, userRequestBody)

	// Create a request
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBufferString(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Set the API key in the request headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read and display the response
	responseBody := new(bytes.Buffer)
	responseBody.ReadFrom(resp.Body)
	//fmt.Println("Response:", responseBody.String())

	responseString := responseBody.String()
	return &responseString, nil
}
