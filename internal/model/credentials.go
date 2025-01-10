// Package model defines the data structures used in the application.
//
// This package includes the Credentials struct, which represents a user's login credentials.
package model

// Credentials represents the login credentials associated with a user.
//
// Fields:
//   - Login: A string representing the user's login name or email address.
//   - Password: A string containing the user's password, which should be stored securely.
//   - Description: A string providing additional information about the credentials.
//   - UserID: An int64 representing the unique identifier of the user associated with these credentials.
//   - ID: An int64 representing the unique identifier of the credentials themselves.
type Credentials struct {
	Login       string
	Password    string
	Description string
	UserID      int64
	ID          int64
}
