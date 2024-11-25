// Package model defines the data structures used in the application.
//
// This package includes the Note struct, which represents a note created by a user.
package model

// Note represents a note associated with a user.
//
// Fields:
//   - Text: A string containing the content of the note.
//   - Description: A string providing additional information about the note.
//   - UserID: An int64 representing the unique identifier of the user who created the note.
//   - ID: An int64 representing the unique identifier of the note itself.
type Note struct {
	Text        string
	Description string
	UserID      int64
	ID          int64
}
