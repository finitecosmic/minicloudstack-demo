package objectstore

import (
	"time"
)

type Object struct {
	Bucket      string
	Key         string
	Size        int64
	ContentType string
	CreatedAt   time.Time
}
