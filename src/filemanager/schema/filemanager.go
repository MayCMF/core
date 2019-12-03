package schema

import (
	"github.com/MayCMF/core/src/common/schema"
)

// file - file object
type file struct {
	UUID     string `json:"uuid"`                        // UUID
	UID      uint   `json:"uid" binding:"required"`      // User ID
	Filename string `json:"filename" binding:"required"` // File Name
	Uri      string `json:"uri"`                         // File URI
	Filemime string `json:"filemime"`                    // Filemime (image/jpeg, application/msword etc)
	Filesize int    `json:"filesize"`                    // Filesize in bytes
}

// fileQueryParam - Query conditions
type fileQueryParam struct {
	UUID         string // UUID
	Filename     uint64 // File Name
	Uri          string // File URI
	LikeFilename string // Name (fuzzy query)
}

// fileQueryOptions - file object query optional parameter item
type fileQueryOptions struct {
	PageParam *schema.PaginationParam // Paging parameter
}

// fileQueryResult - file object query result
type fileQueryResult struct {
	Data       []*file
	PageResult *schema.PaginationResult
}
