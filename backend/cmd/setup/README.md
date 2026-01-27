# Master Reset & Seed Script

## Overview
Complete database reset and population script that wipes everything clean and rebuilds with verified, high-quality data.

## What It Does

1. **🗑️ Hard Reset**
   - Drops ALL existing tables (trims, models, brands)
   - Ensures complete clean slate

2. **🏗️ Schema Recreation**
   - Creates fresh tables with proper foreign keys
   - Includes `image_url` column in trims table
   - Enforces relational integrity

3. **📥 Smart Data Ingestion**
   - Fetches from API Ninjas (BMW, Audi, VW, Fiat)
   - Years: 2023, 2024
   - Maintains proper Brand → Model → Trim hierarchy
   - Prevents duplicates

4. **🖼️ Context-Aware Image Search**
   - **THE FIX**: Uses specific query format
   - Query: `"{Brand} {Model} {Trim} {Year} exterior car studio lighting white background"`
   - Validates image URLs (must end with .jpg/.png/.jpeg)
   - Tries multiple results if first isn't valid
   - Rate limited (1.5s between image searches)

## Usage

```bash
# Make sure your .env has NINJAS_API_KEY
cd c:\Projects\car-specs\backend

# Run the setup
go run cmd/setup/main.go
```

## What You'll See

```
=== Master Database Reset & Seed ===

🗑️  Dropping all tables...
  ✓ Dropped table: trim_features
  ✓ Dropped table: features
  ✓ Dropped table: trims
  ✓ Dropped table: models
  ✓ Dropped table: brands

🏗️  Creating fresh schema...
  ✓ Schema created successfully

📥 Fetching BMW vehicles for year 2023...
  Found 1 vehicles
  ✓ Created brand: bmw
  ✓ Created model: m850i xdrive coupe
  🖼️  Searching for image: bmw m850i xdrive coupe 2023...
  ✓ Found image: https://...
  ✓ Created trim: bmw m850i xdrive coupe 2023

...

📊 Setup Summary:
  Brands created: 4
  Models created: 8
  Trims created:  8
  Images found:   8
  Errors:         0

✅ Setup complete! Database is ready.
```

## Features

✅ **Idempotent** - Safe to run multiple times  
✅ **Context-Aware Images** - No more rock climbers!  
✅ **Image Validation** - Only valid image URLs  
✅ **Rate Limited** - Respects API quotas  
✅ **Error Handling** - Continues on failures  
✅ **Relational Integrity** - Proper foreign keys  

## Troubleshooting

**"NINJAS_API_KEY is required"**
- Add your API key to `.env` file

**"Google API returned 403"**
- Check your Google API key
- Verify Custom Search API is enabled
- Check daily quota (100 free searches/day)

**Images still wrong?**
- The new query format is much more specific
- It includes "exterior car studio lighting" to force car results
- URLs are validated to ensure they're actual images

## After Running

1. Restart your API server:
   ```bash
   .\api_new.exe
   ```

2. Check your frontend:
   - Go to http://localhost:5174/search
   - You should see clean, properly enriched data

3. Verify images:
   - All vehicles should have relevant car images
   - No more random search results
