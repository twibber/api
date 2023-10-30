package lib

var BlankSuccess = Response{
	Success: true,
}

type Response struct {
	Success bool `json:"success"`
	// ObjectName is a pointer to the data just so we're quirky
	ObjectName string      `json:"object,omitempty"`
	Data       any         `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page         int `json:"page"`
	PerPage      int `json:"per_page"`
	PreviousPage int `json:"previous_page"`
	NextPage     int `json:"next_page"`
	LastPage     int `json:"last_page"`
	TotalEntries int `json:"total_entries"`
}
