package dto

import "time"

// NFTResponse represents the response structure for NFT data
type NFTResponse struct {
	TokenID   uint      `json:"token_id" example:"1"`
	Owner     string    `json:"owner" example:"0x1234567890123456789012345678901234567890"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T12:00:00Z"`
}

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"An error occurred"`
	Error   string `json:"error,omitempty" example:"Detailed error message"`
}

// NFTListResponse represents a list of NFTs
type NFTListResponse struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"NFTs retrieved successfully"`
	Data    []NFTResponse `json:"data"`
	Count   int           `json:"count" example:"10"`
}
