package models

// This is our generic response object. This is the format that the portal frontend is expecting so try to stick with this.
type Response[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Status  string `json:"status"`
}
