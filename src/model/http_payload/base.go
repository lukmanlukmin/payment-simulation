// Package httppayload ...
package httppayload

// BaseResponse ...
type BaseResponse struct {
	Success bool          `json:"success"`
	Data    interface{}   `json:"data,omitempty"`
	Meta    *MetaResponse `json:"meta,omitempty"`
	Message *string       `json:"message,omitempty" default:"berhasil"`
	Error   interface{}   `json:"error,omitempty"`
}

// MetaResponse ...
type MetaResponse struct {
	CurrentPage  uint64 `json:"current_page" default:"1"`
	NextPage     *int64 `json:"next_page" default:"2"`
	PreviousPage *int64 `json:"previous_page,omitempty" default:"0"`
	Limit        uint64 `json:"limit" default:"100"`
	TotalData    uint64 `json:"total_data" default:"100"`
	LastPage     uint64 `json:"last_page,omitempty" default:"1"`
}
