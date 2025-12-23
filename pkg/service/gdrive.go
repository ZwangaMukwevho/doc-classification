package service

import (
	"bytes"
	"doc-classification/pkg/model"
	"encoding/base64"
	"fmt"
	"strings"

	"doc-classification/pkg/common"

	"google.golang.org/api/drive/v3"
)

type DriveMethods interface {
	CreateDriveDirectory(name string) (*drive.File, error)
	GetRootDirectoryIDByName(name string) (bool, error)
	FileExists(fileName string, directoryID string) error
	GetDriveDirectories() (*[]model.Directory, error)
	ListFiles(size int64) *[]model.File
	UploadFile(message model.Message) error
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
		common.Logger.Errorf("Unable to retrieve folders: %v", err)
		return nil, err
	}

	// Append them to the directories object
	for _, folder := range driveFiles.Files {
		directory.ID = folder.Id
		directory.Name = folder.Name
		directories = append(directories, directory)
	}

	return &directories, nil
}

func (ds DriveServiceLocal) ListFiles(size int64) *[]model.File {
	var files []model.File
	var file model.File

	results, err := ds.Service.Files.List().PageSize(size).Fields("files(id, name)").Do()
	if err != nil {
		common.Logger.Errorf("Unable to retrieve files: %v", err)
	}

	if len(results.Files) == 0 {
		return nil
	}

	for _, f := range results.Files {
		file.ID = f.Id
		file.Name = f.Name
		files = append(files, file)
	}

	return &files
}

func escapeDriveQueryString(s string) string {
	return strings.ReplaceAll(s, `'`, `\'`)
}

// FileExists checks whether a file with `fileName` exists inside `directoryID`.
// Returns:
// - nil if the file does NOT exist (safe to upload)
// - error if the file already exists (or if the API call fails)
func (ds DriveServiceLocal) FileExists(fileName string, directoryID string) error {
	if fileName == "" {
		return fmt.Errorf("fileName cannot be empty")
	}
	if directoryID == "" {
		return fmt.Errorf("directoryID cannot be empty")
	}

	name := escapeDriveQueryString(fileName)

	// Google Drive v3 query: match by name + parent folder + not trashed
	q := fmt.Sprintf("name = '%s' and '%s' in parents and trashed = false", name, directoryID)

	results, err := ds.Service.Files.
		List().
		Q(q).
		PageSize(1).
		Fields("files(id, name)").
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).
		Do()
	if err != nil {
		common.Logger.Errorf("Unable to query file existence (name=%s, dir=%s): %v", fileName, directoryID, err)
		return err
	}

	if len(results.Files) > 0 {
		f := results.Files[0]
		return fmt.Errorf("file already exists in directory: name=%q dir=%s (fileId=%s)", fileName, directoryID, f.Id)
	}

	return nil
}

func (ds DriveServiceLocal) UploadFile(attachment model.Attachment, directoryID string) error {
	// create file object
	file := &drive.File{
		Name:     attachment.Name,
		MimeType: attachment.MimeType,
		Parents:  []string{directoryID},
	}

	// Decode the base64url data
	data, err := base64.URLEncoding.DecodeString(attachment.Bytestream)
	if err != nil {
		common.Logger.Errorf("Error decoding base64url data: %v", err)
		return err
	}

	// Upload the file
	_, err = ds.Service.Files.Create(file).Media(bytes.NewReader((data))).Do()
	if err != nil {
		common.Logger.Errorf("Unable to create file: %v", err)
		return err
	}

	common.Logger.Infof("File %s uploaded to the specified directory (%s) in Google Drive.", attachment.Name, directoryID)

	return nil
}

func (ds DriveServiceLocal) CreateDriveDirectory(name string) (*drive.File, error) {
	file, err := ds.Service.Files.Create(&drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
	}).Do()

	if err != nil {
		common.Logger.Errorf("Error creating drive categories : %v", err)
		return nil, err
	}

	return file, err
}

func (ds DriveServiceLocal) GetRootDirectoryIDByName(name string) (*string, error) {
	if name == "" {
		common.Logger.Info("Directory name is required when getting")
		return nil, fmt.Errorf("directory name cannot be empty")
	}

	folderName := escapeDriveQueryString(name)

	q := fmt.Sprintf(
		"mimeType = 'application/vnd.google-apps.folder' and name = '%s' and 'root' in parents and trashed = false",
		folderName,
	)

	results, err := ds.Service.Files.
		List().
		Q(q).
		PageSize(1).
		Fields("files(id)").
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).
		Do()

	if err != nil {
		common.Logger.Errorf("Unable to find directory in root (name=%s): %v", name, err)
		return nil, err
	}

	if len(results.Files) == 0 {
		return nil, nil
	}

	fileName := results.Files[0].Id

	return &fileName, nil
}
