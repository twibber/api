package lib

// BlankSuccess provides a default successful response when no additional data is required.
var BlankSuccess = Response{
	Success: true,
}

// Response represents the standard structure for API responses.
type Response struct {
	Success    bool        `json:"success"`              // Indicates if the request was successful
	Data       any         `json:"data,omitempty"`       // Holds the data payload of the response, if any
	Pagination *Pagination `json:"pagination,omitempty"` // Optional pagination details, included for list responses
}

// Pagination details the structure for pagination metadata in list responses.
type Pagination struct {
	Page         int `json:"page"`          // The current page number
	PerPage      int `json:"per_page"`      // The number of items per page
	PreviousPage int `json:"previous_page"` // The previous page number, if applicable
	NextPage     int `json:"next_page"`     // The next page number, if applicable
	LastPage     int `json:"last_page"`     // The last page number based on total entries
	TotalEntries int `json:"total_entries"` // The total number of entries across all pages
}
