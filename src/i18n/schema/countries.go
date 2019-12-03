package schema

import (
	"encoding/json"

	"github.com/MayCMF/core/src/common/schema"
)

// Country - Country object
type Country struct {
	Code      string          `json:"code" binding:"required"` // Country's Code
	Name      string          `json:"name" binding:"required"` // Native
	Native    string          `json:"native"`                  // Native language
	Phone     string          `json:"phone"`                   // Phone number
	Continent string          `json:"continent"`               // Continent where country located
	Capital   string          `json:"capital"`                 // Country's Capital
	Currency  json.RawMessage `json:"currency"`                // Country's Currency
	Languages json.RawMessage `json:"languages"`               // Used languages
	Timezones json.RawMessage `json:"timezones"`               // Country's timezones
	LatLng    json.RawMessage `json:"latlng"`                  // Country's geo center
	Emoji     string          `json:"emoji"`                   // Emoji flag
	EmojiU    string          `json:"emojiU"`                  // Emoji flag
}

// CountryQueryParam - Query conditions
type CountryQueryParam struct {
	Code        string // CODE
	Name        string
	Native      string
	Phone       string
	Continent   string // Continent
	Capital     string // Capital
	Currency    string
	Languages   string // Languages
	LikeCode    string // Number (fuzzy query)
	LikeName    string // Name (fuzzy query)
	LikeCapital string // Name (fuzzy query)
	LikeNative  string // Name (fuzzy query)
}

// CountryQueryOptions - Countries object query optional parameter item
type CountryQueryOptions struct {
	PageParam *schema.PaginationParam // Paging parameter
}

// Countries - Country list
type Countries []*Country

// CountryQueryResult - Countries object query result
type CountryQueryResult struct {
	Data       []*Country
	PageResult *schema.PaginationResult
}
