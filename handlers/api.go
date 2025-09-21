package handlers

import (
	"net/http"
	"strconv"

	"go-cli-eth/dto"
	"go-cli-eth/models"
	"go-cli-eth/services"

	"github.com/gin-gonic/gin"
)

// NFTHandler handles NFT-related API endpoints
type NFTHandler struct {
	nftService *services.NFTService
}

// NewNFTHandler creates a new NFT handler
func NewNFTHandler(nftService *services.NFTService) *NFTHandler {
	return &NFTHandler{
		nftService: nftService,
	}
}

// GetAndStoreOwner godoc
// @Summary Get and store NFT owner from blockchain
// @Description Fetches the owner of an NFT from the blockchain and stores it in the database
// @Tags NFT
// @Accept json
// @Produce json
// @Param request body dto.GetOwnerRequest true "Get owner request"
// @Success 200 {object} dto.SuccessResponse{data=dto.NFTResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/nft/owner [post]
func (h *NFTHandler) GetAndStoreOwner(c *gin.Context) {
	var req dto.GetOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	nft, err := h.nftService.GetAndStoreOwner(req.ContractAddress, req.TokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to get and store NFT owner",
			Error:   err.Error(),
		})
		return
	}

	response := dto.NFTResponse{
		TokenID:   nft.TokenID,
		Owner:     nft.Owner,
		CreatedAt: nft.CreatedAt,
		UpdatedAt: nft.UpdatedAt,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "NFT owner retrieved and stored successfully",
		Data:    response,
	})
}

// UpdateOwner godoc
// @Summary Update NFT owner data
// @Description Updates an existing NFT record with current blockchain data
// @Tags NFT
// @Accept json
// @Produce json
// @Param request body dto.UpdateOwnerRequest true "Update owner request"
// @Success 200 {object} dto.SuccessResponse{data=dto.NFTResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/nft/owner [put]
func (h *NFTHandler) UpdateOwner(c *gin.Context) {
	var req dto.UpdateOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	nft, err := h.nftService.UpdateOwner(req.ContractAddress, req.TokenID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "NFT with token ID "+strconv.Itoa(int(req.TokenID))+" not found in database" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: "Failed to update NFT owner",
			Error:   err.Error(),
		})
		return
	}

	response := dto.NFTResponse{
		TokenID:   nft.TokenID,
		Owner:     nft.Owner,
		CreatedAt: nft.CreatedAt,
		UpdatedAt: nft.UpdatedAt,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "NFT owner updated successfully",
		Data:    response,
	})
}

// GetNFTByTokenID godoc
// @Summary Get NFT by token ID
// @Description Retrieves an NFT record from the database by token ID
// @Tags NFT
// @Produce json
// @Param token_id path int true "Token ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.NFTResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/nft/{token_id} [get]
func (h *NFTHandler) GetNFTByTokenID(c *gin.Context) {
	tokenIDStr := c.Param("token_id")
	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid token ID",
			Error:   err.Error(),
		})
		return
	}

	nft, err := h.nftService.GetNFTByTokenID(uint(tokenID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "NFT with token ID "+tokenIDStr+" not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: "Failed to get NFT",
			Error:   err.Error(),
		})
		return
	}

	response := dto.NFTResponse{
		TokenID:   nft.TokenID,
		Owner:     nft.Owner,
		CreatedAt: nft.CreatedAt,
		UpdatedAt: nft.UpdatedAt,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "NFT retrieved successfully",
		Data:    response,
	})
}

// GetAllNFTs godoc
// @Summary Get all NFTs
// @Description Retrieves all NFT records from the database
// @Tags NFT
// @Produce json
// @Success 200 {object} dto.NFTListResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/nft [get]
func (h *NFTHandler) GetAllNFTs(c *gin.Context) {
	nfts, err := h.nftService.GetAllNFTs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to get NFTs",
			Error:   err.Error(),
		})
		return
	}

	var nftResponses []dto.NFTResponse
	for _, nft := range nfts {
		nftResponses = append(nftResponses, dto.NFTResponse{
			TokenID:   nft.TokenID,
			Owner:     nft.Owner,
			CreatedAt: nft.CreatedAt,
			UpdatedAt: nft.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.NFTListResponse{
		Success: true,
		Message: "NFTs retrieved successfully",
		Data:    nftResponses,
		Count:   len(nftResponses),
	})
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns the health status of the API
// @Tags Health
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Router /health [get]
func (h *NFTHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "API is healthy",
		Data: gin.H{
			"status":  "ok",
			"service": "go-cli-eth",
		},
	})
}

// ConvertModelToDTO converts models.NFT to dto.NFTResponse
func ConvertModelToDTO(nft *models.NFT) dto.NFTResponse {
	return dto.NFTResponse{
		TokenID:   nft.TokenID,
		Owner:     nft.Owner,
		CreatedAt: nft.CreatedAt,
		UpdatedAt: nft.UpdatedAt,
	}
}
