package service

import (
	"bytes"
	"doc-classification/pkg/model"
	"encoding/base64"
	"fmt"
	"log"

	"google.golang.org/api/drive/v3"
)

type DriveMethods interface {
	ListFiles(size int64) *[]model.File
	UploadFile(message model.Message) error
	GetDriveDirectories() (*[]model.Directory, error)
}

type DriveServiceLocal struct {
	Service *drive.Service
}

func (ds DriveServiceLocal) GetDriveDirectories() (*[]model.Directory, error) {
	var directory model.Directory
	var directories []model.Directory

	//Fetch all the directories on the drive
	driveFiles, err := ds.Service.Files.List().Q("mimeType='application/vnd.google-apps.folder'").Do()
	if err != nil {
		log.Printf("Unable to retrieve folders: %v", err)
		return nil, err
	}

	// Append them to the directories object
	for _, folder := range driveFiles.Files {
		directory.ID = folder.Id
		directory.Name = folder.Name
		directories = append(directories, directory)
		fmt.Printf("Name: %s, ID: %s\n", folder.Name, folder.Id)
	}

	return &directories, nil
}

func (ds DriveServiceLocal) ListFiles(size int64) *[]model.File {
	var files []model.File
	var file model.File

	results, err := ds.Service.Files.List().PageSize(size).Fields("files(id, name)").Do()
	if err != nil {
		log.Printf("Unable to retrieve files: %v", err)
	}

	if len(results.Files) == 0 {
		log.Println("No files found on google drive")
	} else {
		for _, f := range results.Files {
			file.ID = f.Id
			file.Name = f.Name
			files = append(files, file)
			fmt.Printf("%s (%s)\n", f.Name, f.Id)
		}
	}
	return &files
}

func (ds DriveServiceLocal) UploadFile(message model.Message, directoryID string) error {
	// create file object
	file := &drive.File{
		Name:     message.File.Name,
		MimeType: message.File.MimeType,
		Parents:  []string{directoryID},
	}

	// Decode the base64url data
	data, err := base64.URLEncoding.DecodeString(message.File.Bytestream)
	if err != nil {
		log.Printf("Error decoding base64url data: %v", err)
		return err
	}

	// Upload the file
	_, err = ds.Service.Files.Create(file).Media(bytes.NewReader((data))).Do()
	if err != nil {
		log.Printf("Unable to create file: %v", err)
		return err
	}
	fmt.Printf("File '%s' uploaded to the specified directory in Google Drive.\n", message.File.Name)
	return nil
}
