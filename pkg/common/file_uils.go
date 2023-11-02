package common

import (
	"fmt"

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
