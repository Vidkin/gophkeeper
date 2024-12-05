// Package model defines the data structures used in the application.
//
// This package includes the User struct, which represents a user in the system.
package model

// User represents a user in the application.
//
// Fields:
//   - Login: A string representing the user's login name or email address.
//   - Password: A string containing the user's password, which should be stored securely.
//   - ID: An int64 representing the unique identifier of the user.
type User struct {
	Login    string
	Password string
	ID       int64
}
