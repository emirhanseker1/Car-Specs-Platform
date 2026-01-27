# Vehicle Ingestion Service - Setup Guide

## Overview
Automated service that fetches vehicle specifications from API Ninjas and enriches them with images from Google Custom Search.

## Prerequisites

### 1. API Ninjas API Key (Required)
1. Go to https://api-ninjas.com/
2. Click "Sign Up" (free tier available)
3. After signup, go to "My Account" → "API Key"
4. Copy your API key
5. Free tier: 50,000 requests/month

### 2. Google Custom Search API (Optional - for images)
**Step A: Get Google API Key**
1. Go to https://console.cloud.google.com/
2. Create a new project or select existing
3. Enable "Custom Search API"
4. Go to "Credentials" → "Create Credentials" → "API Key"
5. Copy your API key
6. Free tier: 100 searches/day

**Step B: Create Custom Search Engine**
1. Go to https://programmablesearchengine.google.com/
2. Click "Add" to create new search engine
3. In "Sites to search": Enter `*` (search entire web)
4. Enable "Image search"
5. Click "Create"
6. Copy your "Search engine ID"

## Installation

1. **Install dependencies:**
```bash
go get github.com/joho/godotenv
```

2. **Create `.env` file:**
```bash
cp .env.example .env
```

3. **Edit `.env` and add your API keys:**
```
NINJAS_API_KEY=your_actual_key_here
GOOGLE_API_KEY=your_google_key_here  # Optional
SEARCH_ENGINE_ID=your_search_id_here  # Optional
```

## Usage

**Run the ingestion service:**
```bash
go run cmd/ingestion/main.go
```

**What it does:**
1. Fetches vehicles from API Ninjas for brands: BMW, Audi, VW, Mercedes, Toyota, Ford
2. For years: 2023, 2024
3. Creates brands and models if they don't exist
4. Checks for duplicates (skips if year already exists)
5. Optionally finds images via Google Custom Search
6. Rate limits: 2 seconds between API calls

## Features

✅ **Idempotent** - Won't create duplicates  
✅ **Rate Limited** - Respects API limits  
✅ **Error Handling** - Continues on errors  
✅ **Relational Integrity** - Proper brand → model → trim hierarchy  
✅ **Data Mapping** - Converts API Ninjas format to your schema  

## Customization

**Change target brands:**
```go
targetBrands := []string{"bmw", "audi", "volkswagen", "mercedes-benz", "toyota", "ford"}
```

**Change target years:**
```go
targetYears := []int{2023, 2024}
```

**Adjust rate limiting:**
```go
time.Sleep(2 * time.Second) // Change to your needs
```

## Troubleshooting

**"NINJAS_API_KEY is required"**
- Make sure `.env` file exists and contains your API key

**"API Ninjas returned status 401"**
- Check your API key is correct
- Verify your API key is active on api-ninjas.com

**"Google API returned status 403"**
- Check your Google API key
- Verify Custom Search API is enabled
- Check you haven't exceeded daily quota (100/day free)

**"Trim already exists for year X, skipping"**
- This is normal - the service is idempotent and won't create duplicates

## Cost

- **API Ninjas**: Free tier = 50,000 requests/month
- **Google Custom Search**: Free tier = 100 searches/day
- **Total Cost**: $0 for moderate usage

## Next Steps

After running ingestion:
1. Check your database: `sqlite3 vehicles.db "SELECT COUNT(*) FROM trims;"`
2. Restart your API server to see new data
3. Visit frontend to browse imported vehicles
