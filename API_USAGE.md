# Go CLI Ethereum NFT Tracker - API Usage Guide

## 🚀 Quick Start

Your Go CLI Ethereum NFT Tracker now supports both CLI and REST API modes with full Swagger documentation!

### Build and Run

```bash
# Build the project
go build -o build/nft-tracker.exe .

# Run in CLI mode (default)
./build/nft-tracker.exe

# Run in API mode
./build/nft-tracker.exe -mode=api -port=8080
```

## 🌐 REST API Endpoints

When running in API mode, the following endpoints are available:

### Health Check
```
GET /health
```

### NFT Operations
```
POST /api/nft/owner          - Get and store NFT owner from blockchain
PUT  /api/nft/owner          - Update existing NFT owner data
GET  /api/nft/{token_id}     - Get NFT by token ID
GET  /api/nft                - List all NFTs
```

### Swagger Documentation
```
GET /swagger/index.html      - Interactive API documentation
```

## 📖 API Examples

### 1. Get and Store NFT Owner
```bash
curl -X POST http://localhost:8080/api/nft/owner \
  -H "Content-Type: application/json" \
  -d '{
    "contract_address": "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D",
    "token_id": 1
  }'
```

### 2. Update NFT Owner
```bash
curl -X PUT http://localhost:8080/api/nft/owner \
  -H "Content-Type: application/json" \
  -d '{
    "contract_address": "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D", 
    "token_id": 1
  }'
```

### 3. Get NFT by Token ID
```bash
curl http://localhost:8080/api/nft/1
```

### 4. List All NFTs
```bash
curl http://localhost:8080/api/nft
```

### 5. Health Check
```bash
curl http://localhost:8080/health
```

## 🔧 Environment Configuration

Set these environment variables or use the `.env` file:

```env
# Database Configuration
DATABASE_URL=postgresql://username:password@localhost:5432/nft_tracker?sslmode=disable

# Ethereum RPC Configuration  
ETH_RPC_URL=https://necessary-evocative-star.bsc-testnet.quiknode.pro/76dbcf8d993da60cf34d714d7eb0e1656c1cf7d7
```

## 📋 Response Format

All API responses follow this format:

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    "token_id": 1,
    "owner": "0x1234567890123456789012345678901234567890",
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z"
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "An error occurred",
  "error": "Detailed error message"
}
```

## 🧪 Testing the API

1. **Start the API server:**
   ```bash
   ./build/nft-tracker.exe -mode=api
   ```

2. **Open Swagger UI:**
   Visit: http://localhost:8080/swagger/index.html

3. **Test endpoints:**
   Use the interactive Swagger UI to test all endpoints with sample data.

## 📁 Project Structure

```
go-cli-eth/
├── main.go                 # Main application with dual CLI/API mode
├── handlers/api.go         # REST API handlers with Swagger annotations
├── dto/                    # Data Transfer Objects
│   ├── requests.go         # API request structures
│   └── responses.go        # API response structures
├── models/nft.go          # Database models
├── database/db.go         # Database connection and setup
├── ethereum/client.go     # Blockchain interaction
├── services/nft_service.go # Business logic layer
├── docs/                  # Auto-generated Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
└── .env                   # Environment configuration
```

## 🎯 Key Features Implemented

✅ **Dual Mode Support**: CLI and REST API modes  
✅ **Swagger Documentation**: Interactive API documentation  
✅ **Environment Configuration**: `.env` file support  
✅ **Error Handling**: Comprehensive error responses  
✅ **CORS Support**: Cross-origin resource sharing  
✅ **Health Check**: Service health monitoring  
✅ **Database Integration**: PostgreSQL with GORM  
✅ **Blockchain Integration**: Ethereum via go-ethereum  

Your NFT tracker is now ready for both development and production use! 🎉