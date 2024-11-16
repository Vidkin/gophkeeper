package model

import "time"

type File struct {
	CreatedAt   time.Time
	BucketName  string
	FileName    string
	ContentType string
	Description string
	UserID      int64
	ID          int64
	FileSize    int64
}
