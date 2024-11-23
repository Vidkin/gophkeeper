package model

type File struct {
	CreatedAt   string
	BucketName  string
	FileName    string
	Description string
	UserID      int64
	ID          int64
	FileSize    int64
}
