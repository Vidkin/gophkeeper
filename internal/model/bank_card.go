// Package model defines the data structures used in the application.
//
// This package includes the BankCard struct, which represents a user's bank card information.
package model

// BankCard represents a bank card associated with a user.
//
// Fields:
//   - ExpireDate: A string representing the expiration date of the bank card (e.g., "MM/YY").
//   - Owner: A string containing the name of the card owner.
//   - CVV: A string representing the card verification value, typically a 3 or 4 digit number.
//   - Number: A string representing the card number, usually 16 digits long.
//   - Description: A string providing additional information about the bank card.
//   - UserID: An int64 representing the unique identifier of the user associated with the bank card.
//   - ID: An int64 representing the unique identifier of the bank card itself.
type BankCard struct {
	ExpireDate  string
	Owner       string
	CVV         string
	Number      string
	Description string
	UserID      int64
	ID          int64
}
