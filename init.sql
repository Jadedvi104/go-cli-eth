-- Initial database setup for NFT Tracker
-- This file is executed when the PostgreSQL container starts for the first time

-- Create the nfts table (GORM will also create it automatically, but this is for reference)
CREATE TABLE IF NOT EXISTS nfts (
    token_id BIGINT PRIMARY KEY,
    owner VARCHAR(42) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on the owner field for faster queries
CREATE INDEX IF NOT EXISTS idx_nfts_owner ON nfts(owner);

-- Create an index on the updated_at field for sorting
CREATE INDEX IF NOT EXISTS idx_nfts_updated_at ON nfts(updated_at);

-- Insert some sample data (optional)
-- INSERT INTO nfts (token_id, owner) VALUES 
-- (1, '0x1234567890123456789012345678901234567890'),
-- (2, '0x0987654321098765432109876543210987654321')
-- ON CONFLICT (token_id) DO NOTHING;
