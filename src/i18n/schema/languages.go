package schema

import (
	"github.com/MayCMF/core/src/common/schema"
)

// Languages - Languages object
type Language struct {
	Code    string `json:"code"`                    // Code
	Name    string `json:"name" binding:"required"` // Name
	Native  string `json:"native"`                  // Remarks
	Rtl     bool   `json:"rtl"`                     // RTL
	Default bool   `json:"default"`                 // Default language for project
	Active  bool   `json:"active"`                  // Status language can be added as content
}

// LanguagesQueryParam - Query conditions
type LanguageQueryParam struct {
	Code       string // CODE
	Status     int    // Status (1: Enable 0: Disable)
	LikeCode   string // Number (fuzzy query)
	LikeName   string // Name (fuzzy query)
	LikeNative string // Name (fuzzy query)
}

// LanguagesQueryOptions - Languages object query optional parameter item
type LanguageQueryOptions struct {
	PageParam *schema.PaginationParam // Paging parameter
}

// Languages - Language list
type Languages []*Language

// LanguagesQueryResult - Languages object query result
type LanguageQueryResult struct {
	Data       []*Language
	PageResult *schema.PaginationResult
}
