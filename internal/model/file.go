// Package model defines the data structures used in the application.
//
// This package includes the File struct, which represents a file associated with a user.
package model

// File represents a file uploaded by a user.
//
// Fields:
//   - CreatedAt: A string representing the date and time when the file was created (e.g., in ISO 8601 format).
//   - BucketName: A string representing the name of the storage bucket where the file is stored.
//   - FileName: A string containing the name of the file, including its extension.
//   - Description: A string providing additional information about the file.
//   - UserID: An int64 representing the unique identifier of the user who uploaded the file.
//   - ID: An int64 representing the unique identifier of the file itself.
//   - FileSize: An int64 representing the size of the file in bytes.
type File struct {
	CreatedAt   string
	BucketName  string
	FileName    string
	Description string
	UserID      int64
	ID          int64
	FileSize    int64
}
