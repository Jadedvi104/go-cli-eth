package services

import (
	"fmt"
	"log"
	"time"

	"go-cli-eth/database"
	"go-cli-eth/ethereum"
	"go-cli-eth/models"

	"gorm.io/gorm"
)

// NFTService handles NFT operations
type NFTService struct {
	ethClient *ethereum.EthereumClient
}

// NewNFTService creates a new NFT service
func NewNFTService(ethClient *ethereum.EthereumClient) *NFTService {
	return &NFTService{
		ethClient: ethClient,
	}
}

// GetAndStoreOwner retrieves owner from blockchain and stores in database
func (s *NFTService) GetAndStoreOwner(contractAddress string, tokenID uint) (*models.NFT, error) {
	// Get owner from blockchain
	owner, err := s.ethClient.GetOwnerOf(contractAddress, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner from blockchain: %v", err)
	}

	log.Printf("Retrieved owner %s for token ID %d", owner, tokenID)

	// Check if NFT already exists in database
	var existingNFT models.NFT
	db := database.GetDB()

	err = db.Where("token_id = ?", tokenID).First(&existingNFT).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing NFT: %v", err)
	}

	// If NFT exists, return existing record
	if err != gorm.ErrRecordNotFound {
		log.Printf("NFT with token ID %d already exists in database", tokenID)
		return &existingNFT, nil
	}

	// Create new NFT record
	nft := models.NFT{
		TokenID:   tokenID,
		Owner:     owner,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	err = db.Create(&nft).Error
	if err != nil {
		return nil, fmt.Errorf("failed to save NFT to database: %v", err)
	}

	log.Printf("Successfully stored NFT with token ID %d", tokenID)
	return &nft, nil
}

// UpdateOwner updates the owner of an existing NFT
func (s *NFTService) UpdateOwner(contractAddress string, tokenID uint) (*models.NFT, error) {
	// Get current owner from blockchain
	owner, err := s.ethClient.GetOwnerOf(contractAddress, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner from blockchain: %v", err)
	}

	log.Printf("Retrieved updated owner %s for token ID %d", owner, tokenID)

	// Update in database
	db := database.GetDB()
	var nft models.NFT

	err = db.Where("token_id = ?", tokenID).First(&nft).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("NFT with token ID %d not found in database", tokenID)
		}
		return nil, fmt.Errorf("failed to find NFT: %v", err)
	}

	// Update the owner and timestamp
	nft.Owner = owner
	nft.UpdatedAt = time.Now()

	err = db.Save(&nft).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update NFT: %v", err)
	}

	log.Printf("Successfully updated NFT with token ID %d", tokenID)
	return &nft, nil
}

// GetNFTByTokenID retrieves an NFT by token ID from database
func (s *NFTService) GetNFTByTokenID(tokenID uint) (*models.NFT, error) {
	db := database.GetDB()
	var nft models.NFT

	err := db.Where("token_id = ?", tokenID).First(&nft).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("NFT with token ID %d not found", tokenID)
		}
		return nil, fmt.Errorf("failed to get NFT: %v", err)
	}

	return &nft, nil
}

// GetAllNFTs retrieves all NFTs from database
func (s *NFTService) GetAllNFTs() ([]models.NFT, error) {
	db := database.GetDB()
	var nfts []models.NFT

	err := db.Find(&nfts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get NFTs: %v", err)
	}

	return nfts, nil
}
