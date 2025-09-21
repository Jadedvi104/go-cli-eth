// Package main provides a REST API server for the Go Ethereum NFT Tracker
// @title Go Ethereum NFT Tracker API
// @version 1.0
// @description This is a REST API server for tracking NFT ownership data from the Ethereum blockchain and storing it in a PostgreSQL database.
//
// @contact.name API Support
// @contact.url http://github.com/Jadedvi104/go-cli-eth
// @contact.email support@example.com
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"go-cli-eth/database"
	_ "go-cli-eth/docs" // Import for swagger docs
	"go-cli-eth/ethereum"
	"go-cli-eth/handlers"
	"go-cli-eth/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	// Parse command line flags
	var (
		mode = flag.String("mode", "api", "Run mode: 'api' for command line interface, 'api' for HTTP API server")
		port = flag.String("port", "8080", "Port for API server (only used in api mode)")
	)
	flag.Parse()

	fmt.Println("🚀 Go CLI Ethereum NFT Tracker")
	fmt.Println("===============================")

	// Initialize database connection
	var dbConnectionString string
	if *mode == "cli" {
		fmt.Print("Enter PostgreSQL connection string (or press Enter to use DATABASE_URL env var): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		dbConnectionString = strings.TrimSpace(input)
	} else {
		// In API mode, try to get from environment
		dbConnectionString = os.Getenv("DATABASE_URL")
		if dbConnectionString == "" {
			log.Fatal("DATABASE_URL environment variable is required in API mode")
		}
		fmt.Printf("Using database from environment: %s\n", maskConnectionString(dbConnectionString))
	}

	err := database.InitDB(dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Ethereum client
	var rpcURL string
	if *mode == "cli" {
		fmt.Print("Enter Ethereum RPC URL (or press Enter for default): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		rpcURL = strings.TrimSpace(input)
	}

	if rpcURL == "" {
		rpcURL = defaultRPCURL
		fmt.Printf("Using RPC URL: %s\n", maskRPCURL(rpcURL))
	}

	ethClient, err := ethereum.NewEthereumClient(rpcURL)
	if err != nil {
		log.Fatalf("Failed to initialize Ethereum client: %v", err)
	}
	defer ethClient.Close()

	// Initialize NFT service
	nftService := services.NewNFTService(ethClient)

	// Run in the specified mode
	switch *mode {
	case "cli":
		runCLI(nftService)
	case "api":
		runAPIServer(nftService, *port)
	default:
		log.Fatalf("Invalid mode: %s. Use 'cli' or 'api'", *mode)
	}
}

// runCLI runs the application in command line interface mode
func runCLI(nftService *services.NFTService) {
	reader := bufio.NewReader(os.Stdin)

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

// runAPIServer runs the application in HTTP API server mode
func runAPIServer(nftService *services.NFTService, port string) {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handlers
	nftHandler := handlers.NewNFTHandler(nftService)

	// Health check endpoint
	r.GET("/health", nftHandler.HealthCheck)

	// API routes
	api := r.Group("/api")
	{
		nft := api.Group("/nft")
		{
			nft.POST("/owner", nftHandler.GetAndStoreOwner)
			nft.PUT("/owner", nftHandler.UpdateOwner)
			nft.GET("/:token_id", nftHandler.GetNFTByTokenID)
			nft.GET("", nftHandler.GetAllNFTs)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Printf("🌐 API server starting on port %s\n", port)
	fmt.Printf("📖 Swagger documentation available at: http://localhost:%s/swagger/index.html\n", port)
	fmt.Printf("🏥 Health check available at: http://localhost:%s/health\n", port)

	// Start server
	log.Fatal(r.Run(":" + port))
}

// maskConnectionString masks sensitive parts of the database connection string
func maskConnectionString(connStr string) string {
	// Simple masking - replace password with ***
	parts := strings.Split(connStr, ":")
	if len(parts) >= 3 {
		// Find password part (after second colon, before @)
		for i, part := range parts {
			if i >= 2 && strings.Contains(part, "@") {
				atIndex := strings.Index(part, "@")
				parts[i] = "***" + part[atIndex:]
				break
			}
		}
		return strings.Join(parts, ":")
	}
	return connStr
}

// maskRPCURL masks sensitive parts of the RPC URL
func maskRPCURL(rpcURL string) string {
	// Simple masking for API keys in URLs
	if strings.Contains(rpcURL, "/") {
		parts := strings.Split(rpcURL, "/")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			if len(lastPart) > 8 {
				parts[len(parts)-1] = lastPart[:4] + "***" + lastPart[len(lastPart)-4:]
			}
		}
		return strings.Join(parts, "/")
	}
	return rpcURL
}

// CLI Handler Functions
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
