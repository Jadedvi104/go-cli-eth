package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthereumClient represents the Ethereum client
type EthereumClient struct {
	client      *ethclient.Client
	contractABI abi.ABI
}

// NewEthereumClient creates a new Ethereum client
func NewEthereumClient(rpcURL string) (*EthereumClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	// ERC-721 ownerOf function ABI
	abiJSON := `[{
		"constant": true,
		"inputs": [{"name": "tokenId", "type": "uint256"}],
		"name": "ownerOf",
		"outputs": [{"name": "", "type": "address"}],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	}]`

	contractABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %v", err)
	}

	return &EthereumClient{
		client:      client,
		contractABI: contractABI,
	}, nil
}

// GetOwnerOf calls the ownerOf function on an ERC-721 contract
func (ec *EthereumClient) GetOwnerOf(contractAddress string, tokenID uint) (string, error) {
	// Convert contract address to common.Address
	contractAddr := common.HexToAddress(contractAddress)

	// Convert tokenID to big.Int
	tokenIDBig := big.NewInt(int64(tokenID))

	// Prepare the call data
	data, err := ec.contractABI.Pack("ownerOf", tokenIDBig)
	if err != nil {
		return "", fmt.Errorf("failed to pack function call: %v", err)
	}

	// Create the call message
	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	// Make the call
	result, err := ec.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", fmt.Errorf("failed to call contract: %v", err)
	}

	// Unpack the result
	var owner common.Address
	err = ec.contractABI.UnpackIntoInterface(&owner, "ownerOf", result)
	if err != nil {
		return "", fmt.Errorf("failed to unpack result: %v", err)
	}

	return owner.Hex(), nil
}

// Close closes the Ethereum client connection
func (ec *EthereumClient) Close() {
	ec.client.Close()
}
