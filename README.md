# Go CLI Ethereum NFT Tracker

A Go CLI application that interacts with Ethereum blockchain to fetch NFT ownership data and store it in a PostgreSQL database.

## Features

1. **Get NFT Owner**: Fetch the owner of an NFT token from the blockchain using the `ownerOf` function
2. **Store Data**: Store retrieved NFT data in PostgreSQL database with tokenId as primary key
3. **Update Function**: Update existing NFT records with current blockchain data
4. **Database Management**: Automatic table creation and data persistence

## Prerequisites

- Go 1.19 or later
- PostgreSQL database
- Ethereum RPC endpoint (Infura, Alchemy, or local node)

## Installation

1. Clone or download this project
2. Navigate to the project directory:
   ```bash
   cd go-cli-eth
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Configuration

1. Copy the example environment file:
   ```bash
   copy .env.example .env
   ```

2. Edit `.env` file with your database and Ethereum RPC configuration:
   ```
   DATABASE_URL=postgres://username:password@localhost:5432/dbname?sslmode=disable
   ETH_RPC_URL=https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID
   ```

### PostgreSQL Setup

Make sure you have PostgreSQL running and create a database:

```sql
CREATE DATABASE nft_tracker;
```

The application will automatically create the required tables.

### Ethereum RPC Setup

You can use:
- **Infura**: Sign up at [infura.io](https://infura.io) and get your project ID
- **Alchemy**: Sign up at [alchemy.com](https://www.alchemy.com) and get your API key
- **Local node**: Run your own Ethereum node

## Usage

1. Build and run the application:
   ```bash
   go run main.go
   ```

2. The application will prompt you to enter:
   - PostgreSQL connection string (or use environment variable)
   - Ethereum RPC URL (or use default from environment)

3. Use the interactive menu to:
   - **Option 1**: Get and store NFT owner data from blockchain
   - **Option 2**: Update existing NFT record with current blockchain data
   - **Option 3**: Retrieve NFT data from database by token ID
   - **Option 4**: List all stored NFTs
   - **Option 5**: Exit

## Project Structure

```
go-cli-eth/
├── main.go                 # Main application entry point
├── models/
│   └── nft.go             # NFT data model
├── database/
│   └── db.go              # Database connection and setup
├── ethereum/
│   └── client.go          # Ethereum client and contract interaction
├── services/
│   └── nft_service.go     # Business logic layer
├── .env.example           # Environment configuration example
├── go.mod                 # Go module file
└── README.md              # This file
```

## Database Schema

The application creates a `nfts` table with the following structure:

| Column    | Type      | Description                 |
|-----------|-----------|-----------------------------|
| token_id  | uint      | Primary key, NFT token ID   |
| owner     | varchar   | Ethereum address of owner   |
| created_at| timestamp | Record creation time        |
| updated_at| timestamp | Last update time            |

## Example Usage

1. **Fetching NFT Owner**: 
   - Enter contract address: `0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D` (BAYC)
   - Enter token ID: `1`
   - Application fetches owner from blockchain and stores in database

2. **Updating NFT Data**:
   - Enter same contract and token ID
   - Application fetches current owner and updates database record

## Error Handling

The application includes comprehensive error handling for:
- Database connection failures
- Ethereum RPC connection issues
- Invalid contract addresses or token IDs
- Network timeouts
- Data validation errors

## Dependencies

- `github.com/ethereum/go-ethereum`: Ethereum client library
- `gorm.io/gorm`: ORM for database operations
- `gorm.io/driver/postgres`: PostgreSQL driver for GORM

## Building for Production

To build a standalone executable:

```bash
go build -o nft-tracker main.go
```

## Security Considerations

- Store sensitive configuration (database passwords, API keys) in environment variables
- Use SSL connections for database and HTTPS for RPC endpoints
- Consider rate limiting when making blockchain calls
- Validate all user inputs

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source and available under the MIT License.
