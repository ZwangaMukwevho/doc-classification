package service

import (
	"bytes"
	"doc-classification/pkg/model"
	"encoding/json"
	"errors"
	"fmt"
)

func CreateClassifyEmailPrompt(message model.Message) string {
	emailSubject := message.Subject
	attachmentName := message.File.Name

	formattedMessage := fmt.Sprintf("From the above categories, how can you classify an attachment from an email with subject '%s' and attachment name: '%s' ", emailSubject, attachmentName)
	return formattedMessage
}

func CreateInitialEmailPrompt() string {
	requestBody := `{
	    "model": "gpt-3.5-turbo",
	    "messages": [
	        {"role": "system", "content": "You are a document classification assistant. I have categories I want to classify my email documents which are: 1. Education: This category can include documents related to educational pursuits, such as school transcripts, certificates, course materials, and research papers. 2. Finance: Finance-related documents can cover a wide range of items, including bank statements, tax records, invoices, quotes, receipts, and investment reports. 5. Work: Work-related documents can involve project plans, reports, emails, resumes, and other materials directly related to your professional life. 6. Home: Home category files may include property documents, utility bills, home maintenance records, and household inventory. 7. Personal: This category can cover a wide range of personal documents, from family photos to personal notes, travel itineraries, and more. Give your reply as one word answer from the given categories"},
			{"role": "user", "content": "How can you classify attachment from an email with subject 'Quotation:21026480' and attachment name: 'Vehicle Booking Acceptance form ref::2102648'"}
	    ]
	}`
	return requestBody
}

func CreateSubsequentPrompt(prompt string) string {
	// Construct the request payload dynamically
	requestBody := fmt.Sprintf(`{"role": "user", "content": "%s"}`, prompt)
	return requestBody
}

func ExtractOpenAIContent(responseBody string) (*string, *error) {
	// Create a bytes buffer and read the response body into it
	responseBuffer := bytes.NewBufferString(responseBody)

	// Decode JSON response
	var jsonResponse map[string]interface{}
	if err := json.NewDecoder(responseBuffer).Decode(&jsonResponse); err != nil {
		return nil, &err
	}

	// Extract the content
	choices, ok := jsonResponse["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		err := errors.New("no choices found in the response")
		return nil, &err
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		err := errors.New("error extracting choice information")
		return nil, &err
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		err := errors.New("error extracting message information")
		return nil, &err
	}

	content, ok := message["content"].(string)
	if !ok {
		err := errors.New("error extracting content")
		return nil, &err
	}

	return &content, nil
}
