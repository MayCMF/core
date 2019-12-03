package schema

import (
	"github.com/MayCMF/core/src/common/schema"
	"time"
)

// file - file object
type File struct {
	UUID      string    `json:"uuid"`                        // UUID
	UID       uint      `json:"uid" binding:"required"`      // User ID
	Filename  string    `json:"filename" binding:"required"` // File Name
	Uri       string    `json:"uri"`                         // File URI
	Filemime  string    `json:"filemime"`                    // Filemime (image/jpeg, application/msword etc)
	Filesize  uint64    `json:"filesize"`                    // Filesize in bytes
	CreatedAt time.Time `json:"created"`                     // File created
}

// fileQueryParam - Query conditions
type FileQueryParam struct {
	UUID         string // UUID
	Filename     string // File Name
	Uri          string // File URI
	LikeFilename string // Name (fuzzy query)
}

// fileQueryOptions - file object query optional parameter item
type FileQueryOptions struct {
	PageParam *schema.PaginationParam // Paging parameter
}

// fileQueryResult - file object query result
type FileQueryResult struct {
	Data       []*File
	PageResult *schema.PaginationResult
}
