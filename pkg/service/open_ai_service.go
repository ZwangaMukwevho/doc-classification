package service

import (
	"doc-classification/pkg/model"
	"fmt"
)

func CreateClassifyEmailPrompt(message model.Message) string {
	emailSubject := message.Subject
	attachmentName := message.File.Name

	formattedMessage := fmt.Sprintf("How can you classify attachment from an email with subject '%s' and attachment name: '%s'", emailSubject, attachmentName)
	return formattedMessage
}

func CreateInitialEmailPromt() string {
	requestBody := `{
        "model": "gpt-3.5-turbo",
        "messages": [
            {"role": "system", "content": "You are a document classification assistant. I have categories I want to classify my email documents which are: 1. Education: This category can include documents related to educational pursuits, such as school transcripts, certificates, course materials, and research papers. 2. Finance: Finance-related documents can cover a wide range of items, including bank statements, tax records, invoices, quotes, receipts, and investment reports. 5. Work: Work-related documents can involve project plans, reports, emails, resumes, and other materials directly related to your professional life. 6. Home: Home category files may include property documents, utility bills, home maintenance records, and household inventory. 7. Personal: This category can cover a wide range of personal documents, from family photos to personal notes, travel itineraries, and more. Give your reply as one word answer from the given categories"},
        ]
    }`
	return requestBody
}

func CreateSubsequentPromts(prompt string) string {
	// Construct the request payload dynamically
	requestBody := fmt.Sprintf(`{
        "model": "gpt-3.5-turbo",
        "messages": [
            {"role": "user", "content": "%s"}
        ]
    }`, prompt)
	return requestBody
}
