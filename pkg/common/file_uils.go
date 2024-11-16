package common

import (
	"doc-classification/pkg/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func GetMimeType(fileName string) string {
	// Create a sample file name with the provided extension
	// fileName := "sample" + fileExtension

	// Get the MIME type using the "mimetype" library
	mime, err := mimetype.DetectFile(fileName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("File name: %s\n", fileName)
		fmt.Printf("MIME type: %s\n", mime.String())
	}

	return mime.String()
}

// FindDirectoryByID searches for a directory by Name and returns its corresponding ID
func FindDirectoryByID(directories map[string]model.Category, nameToFind string) (*string, error) {
	// Iterate through directories to find a match
	defaultDirectoryID := ""
	for ID, dir := range directories {
		if strings.EqualFold(dir.Category, nameToFind) {
			return &ID, nil
		}

		// Get the default ID here
		// This could be set globally somewhere, but we are gonna loop through it here either way, so set here
		if dir.Category == "Default" {
			defaultDirectoryID = ID
		}
	}

	return &defaultDirectoryID, nil
}

func ReadJsonFile(jsonFilePath string) (*[]model.Directory, *error) {
	var directories []model.Directory
	// Get the absolute path to the JSON file
	absolutePath, err := filepath.Abs(jsonFilePath)
	if err != nil {
		log.Println("Error getting absolute path:", err)
		return nil, &err
	}

	// Read JSON data from the file
	jsonData, err := ioutil.ReadFile(absolutePath)
	if err != nil {
		log.Println("Error reading JSON file:", err)
		return nil, &err
	}

	// Unmarshal JSON data into a slice of Directory struct
	err = json.Unmarshal(jsonData, &directories)
	if err != nil {
		log.Println("Error unmarshalling JSON data:", err)
		return nil, &err
	}

	return &directories, nil
}

func GetJsonFileByteStream(filePath string) (*[]byte, error) {
	byteStream, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Unable to read client secret file: %v", err)
		return nil, err
	}
	return &byteStream, nil
}
