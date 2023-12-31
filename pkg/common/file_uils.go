package common

import (
	"doc-classification/pkg/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
func FindDirectoryByID(directories []model.Directory, nameToFind string) (*string, *error) {
	// Iterate through directories to find a match
	defaultDirectoryID := ""
	for _, dir := range directories {
		if strings.EqualFold(dir.Name, nameToFind) {
			return &dir.ID, nil
		}

		// Get the default ID here
		// This could be set globally somewhere, but we are gonna loop through it here either way, so set here
		if dir.Name == "Default" {
			defaultDirectoryID = dir.ID
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
