package model

type BankCard struct {
	ExpireDate  string
	Owner       string
	CVV         string
	Number      string
	Description string
	UserID      int64
	ID          int64
}
