package service

import (
	"doc-classification/pkg/common"
	"doc-classification/pkg/model"

	"google.golang.org/api/gmail/v1"
)

type GmailMethods interface {
	GetAttachmentArray(user string, query string) (*[]model.Message, error)
}

type GmailServiceLocal struct {
	Service *gmail.Service
}

func (gs GmailServiceLocal) GetAttachmentArray(user string, query string) (*[]model.Message, error) {
	messages, err := gs.Service.Users.Messages.List(user).Q(query).Do()
	if err != nil {
		common.Logger.Errorf("Unable to retrieve unread messages: %v", err)
		return nil, err
	}

	// Initialise variables
	var messagesArray []model.Message
	var messageStruct model.Message

	// Loop through all the messages array
	for _, message := range messages.Messages {

		// Get messages
		msg, err := gs.Service.Users.Messages.Get(user, message.Id).Do()
		if err != nil {
			common.Logger.Errorf("Unable to retrieve message details: %v", err)
			continue
		}

		headers := msg.Payload.Headers
		timestamp := ""
		for _, header := range headers {

			if header.Name == "Date" {
				timestamp = header.Value
				messageStruct.Timestamp = timestamp
			}

			if header.Name == "Subject" {
				messageStruct.Subject = header.Value
				break
			}
		}

		// Check for attachments
		parts := msg.Payload.Parts
		var attachments []model.Attachment

		if parts == nil {
			continue
		}

		for _, part := range parts {
			if part.Filename == "" {
				continue
			}

			attachmentID := part.Body.AttachmentId
			if attachmentID == "" {
				continue
			}

			attachmentData, err := gs.Service.Users.Messages.Attachments.Get(user, message.Id, attachmentID).Do()
			if err != nil {
				common.Logger.Errorf("Unable to retrieve attachment content: %v", err)
				continue
			}

			attachment := model.Attachment{
				ID:         attachmentID,
				Name:       part.Filename,
				MimeType:   part.MimeType,
				Bytestream: attachmentData.Data,
				Size:       attachmentData.Size,
			}

			attachments = append(attachments, attachment)
		}

		messageStruct.Files = attachments
		messagesArray = append(messagesArray, messageStruct)
	}
	return &messagesArray, nil
}
