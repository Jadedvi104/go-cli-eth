package models

import (
	"time"
)

// NFT represents the NFT data structure in the database
type NFT struct {
	TokenID   uint      `gorm:"primaryKey" json:"token_id"`
	Owner     string    `gorm:"type:varchar(42);not null" json:"owner"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName returns the table name for the NFT model
func (NFT) TableName() string {
	return "nfts"
}
