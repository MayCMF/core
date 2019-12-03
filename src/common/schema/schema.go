package schema

// HTTPStatusText - Define HTTP status text
type HTTPStatusText string

func (t HTTPStatusText) String() string {
	return string(t)
}

// Define HTTP status text constants
const (
	OKStatusText HTTPStatusText = "OK"
)

// HTTPError - HTTP response error
type HTTPError struct {
	Error HTTPErrorItem `json:"error"` // Error item
}

// HTTPErrorItem HTTP response error item
type HTTPErrorItem struct {
	Code    int    `json:"code"`    // Error code
	Message string `json:"message"` // Error message
}

// HTTPStatus - HTTP response status
type HTTPStatus struct {
	Status string `json:"status"` // status(OK)
}

// HTTPList - HTTP response list data
type HTTPList struct {
	List       interface{}     `json:"list"`
	Pagination *HTTPPagination `json:"pagination,omitempty"`
}

// HTTPPagination - HTTP paging data
type HTTPPagination struct {
	Total    int `json:"total"`
	Current  int `json:"current"`
	PageSize int `json:"pageSize"`
}

// PaginationParam - Paging query condition
type PaginationParam struct {
	PageIndex int // Page index
	PageSize  int // Page size
}

// PaginationResult - Paging query results
type PaginationResult struct {
	Total int // Total number of data
}
