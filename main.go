package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"go-cli-eth/database"
	"go-cli-eth/ethereum"
	"go-cli-eth/services"
)

var (
	// defaultRPCURL is populated from the environment variable ETH_RPC_URL if available.
	// If not present the placeholder Infura URL is used.
	defaultRPCURL string
)

// loadDotEnv loads a simple .env file (KEY=VALUE per line). Comments (#) and empty
// lines are ignored. Existing environment variables are not overwritten.
func loadDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		// don't treat missing .env as an error; just return
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// split at first '='
		if idx := strings.Index(line, "="); idx != -1 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			// remove optional surrounding quotes
			if len(val) >= 2 && ((val[0] == '\'' && val[len(val)-1] == '\'') || (val[0] == '"' && val[len(val)-1] == '"')) {
				val = val[1 : len(val)-1]
			}
			if key != "" {
				if _, exists := os.LookupEnv(key); !exists {
					os.Setenv(key, val)
				}
			}
		}
	}
}

func init() {
	// try to load a local .env file (if present)
	loadDotEnv(".env")

	if v := os.Getenv("ETH_RPC_URL"); v != "" {
		defaultRPCURL = v
	} else {
		defaultRPCURL = "https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID"
	}
}

func main() {
	fmt.Println("🚀 Go CLI Ethereum NFT Tracker")
	fmt.Println("===============================")

	// Initialize database connection
	var dbConnectionString string
	fmt.Print("Enter PostgreSQL connection string (or press Enter to use DATABASE_URL env var): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	dbConnectionString = strings.TrimSpace(input)

	err := database.InitDB(dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Ethereum client
	var rpcURL string
	fmt.Print("Enter Ethereum RPC URL (or press Enter for default): ")
	input, _ = reader.ReadString('\n')
	rpcURL = strings.TrimSpace(input)

	if rpcURL == "" {
		rpcURL = defaultRPCURL
		fmt.Printf("Using default RPC URL: %s\n", rpcURL)
	}

	ethClient, err := ethereum.NewEthereumClient(rpcURL)
	if err != nil {
		log.Fatalf("Failed to initialize Ethereum client: %v", err)
	}
	defer ethClient.Close()

	// Initialize NFT service
	nftService := services.NewNFTService(ethClient)

	// Main application loop
	for {
		fmt.Println("\n📋 Available Commands:")
		fmt.Println("1. Get and store NFT owner")
		fmt.Println("2. Update NFT owner")
		fmt.Println("3. Get NFT by Token ID")
		fmt.Println("4. List all NFTs")
		fmt.Println("5. Exit")
		fmt.Print("\nSelect an option (1-5): ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			handleGetAndStoreOwner(nftService, reader)
		case "2":
			handleUpdateOwner(nftService, reader)
		case "3":
			handleGetNFT(nftService, reader)
		case "4":
			handleListAllNFTs(nftService)
		case "5":
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid option. Please select 1-5.")
		}
	}
}

func handleGetAndStoreOwner(nftService *services.NFTService, reader *bufio.Reader) {
	fmt.Print("Enter contract address: ")
	contractAddress, _ := reader.ReadString('\n')
	contractAddress = strings.TrimSpace(contractAddress)

	fmt.Print("Enter token ID: ")
	tokenIDStr, _ := reader.ReadString('\n')
	tokenIDStr = strings.TrimSpace(tokenIDStr)

	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 32)
	if err != nil {
		fmt.Printf("❌ Invalid token ID: %v\n", err)
		return
	}

	nft, err := nftService.GetAndStoreOwner(contractAddress, uint(tokenID))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Success! NFT Details:\n")
	fmt.Printf("   Token ID: %d\n", nft.TokenID)
	fmt.Printf("   Owner: %s\n", nft.Owner)
	fmt.Printf("   Created: %s\n", nft.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Updated: %s\n", nft.UpdatedAt.Format("2006-01-02 15:04:05"))
}

func handleUpdateOwner(nftService *services.NFTService, reader *bufio.Reader) {
	fmt.Print("Enter contract address: ")
	contractAddress, _ := reader.ReadString('\n')
	contractAddress = strings.TrimSpace(contractAddress)

	fmt.Print("Enter token ID: ")
	tokenIDStr, _ := reader.ReadString('\n')
	tokenIDStr = strings.TrimSpace(tokenIDStr)

	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 32)
	if err != nil {
		fmt.Printf("❌ Invalid token ID: %v\n", err)
		return
	}

	nft, err := nftService.UpdateOwner(contractAddress, uint(tokenID))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Updated! NFT Details:\n")
	fmt.Printf("   Token ID: %d\n", nft.TokenID)
	fmt.Printf("   Owner: %s\n", nft.Owner)
	fmt.Printf("   Created: %s\n", nft.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Updated: %s\n", nft.UpdatedAt.Format("2006-01-02 15:04:05"))
}

func handleGetNFT(nftService *services.NFTService, reader *bufio.Reader) {
	fmt.Print("Enter token ID: ")
	tokenIDStr, _ := reader.ReadString('\n')
	tokenIDStr = strings.TrimSpace(tokenIDStr)

	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 32)
	if err != nil {
		fmt.Printf("❌ Invalid token ID: %v\n", err)
		return
	}

	nft, err := nftService.GetNFTByTokenID(uint(tokenID))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("📄 NFT Details:\n")
	fmt.Printf("   Token ID: %d\n", nft.TokenID)
	fmt.Printf("   Owner: %s\n", nft.Owner)
	fmt.Printf("   Created: %s\n", nft.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Updated: %s\n", nft.UpdatedAt.Format("2006-01-02 15:04:05"))
}

func handleListAllNFTs(nftService *services.NFTService) {
	nfts, err := nftService.GetAllNFTs()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(nfts) == 0 {
		fmt.Println("📭 No NFTs found in database.")
		return
	}

	fmt.Printf("📋 Found %d NFT(s):\n", len(nfts))
	fmt.Println("===========================================")
	for _, nft := range nfts {
		fmt.Printf("Token ID: %d | Owner: %s | Updated: %s\n",
			nft.TokenID,
			nft.Owner,
			nft.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
}
