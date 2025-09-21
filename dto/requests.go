package dto

// GetOwnerRequest represents a request to get NFT owner data
type GetOwnerRequest struct {
	ContractAddress string `json:"contract_address" binding:"required" example:"0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"`
	TokenID         uint   `json:"token_id" binding:"required" example:"1"`
}

// UpdateOwnerRequest represents a request to update NFT owner data
type UpdateOwnerRequest struct {
	ContractAddress string `json:"contract_address" binding:"required" example:"0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"`
	TokenID         uint   `json:"token_id" binding:"required" example:"1"`
}
