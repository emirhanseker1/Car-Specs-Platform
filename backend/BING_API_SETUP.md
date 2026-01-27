# How to Get Bing Image Search API Key

## Overview
Microsoft Bing Image Search API is part of Azure Cognitive Services. It's simpler and more accessible than Google Custom Search, with a generous free tier.

## Free Tier Benefits
- ✅ **1,000 transactions/month FREE**
- ✅ **3 transactions/second**
- ✅ No credit card required for trial
- ✅ Simple setup (no search engine configuration needed)

## Step-by-Step Setup

### 1. Create Azure Account
1. Go to https://azure.microsoft.com/
2. Click "Start free" or "Sign in"
3. Sign in with your Microsoft account (or create one)
4. You get $200 free credit for 30 days (but free tier doesn't need it)

### 2. Create Bing Search Resource

1. **Go to Azure Portal**
   - Visit: https://portal.azure.com/

2. **Create Resource**
   - Click "+ Create a resource"
   - Search for "Bing Search v7"
   - Click "Create"

3. **Configure Resource**
   - **Subscription**: Select your subscription
   - **Resource Group**: Create new or use existing
   - **Region**: Choose closest to you (e.g., "East US")
   - **Name**: `car-specs-image-search` (or any name)
   - **Pricing Tier**: Select **F1 (Free)** - 1K transactions/month

4. **Review + Create**
   - Click "Review + create"
   - Click "Create"
   - Wait for deployment (takes ~1 minute)

### 3. Get Your API Key

1. **Go to Resource**
   - Click "Go to resource" after deployment
   - Or find it in "All resources"

2. **Copy API Key**
   - In the left menu, click "Keys and Endpoint"
   - You'll see two keys (KEY 1 and KEY 2)
   - Copy **KEY 1** (you can use either)
   - It looks like: `a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6`

3. **Add to .env File**
   ```env
   BING_API_KEY=your_actual_key_here
   ```

## Alternative: Quick Trial (No Azure Account)

If you want to test quickly without Azure:

1. Go to https://www.microsoft.com/en-us/bing/apis/bing-image-search-api
2. Click "Try now"
3. Sign in with Microsoft account
4. Get a 7-day trial key (1,000 calls)

## Verify It Works

Test your API key with curl:

```bash
curl -H "Ocp-Apim-Subscription-Key: YOUR_KEY_HERE" \
  "https://api.bing.microsoft.com/v7.0/images/search?q=bmw+320i+2023+car"
```

You should get a JSON response with image results.

## Pricing (After Free Tier)

If you exceed 1,000 searches/month:
- **S1 Tier**: $3 per 1,000 transactions
- Still very affordable for most use cases

## Comparison: Bing vs Google

| Feature | Bing Image Search | Google Custom Search |
|---------|------------------|---------------------|
| Free Tier | 1,000/month | 100/day |
| Setup Complexity | Simple | Complex (need search engine) |
| Credit Card | Not required | Required |
| Image Quality | Excellent | Excellent |
| **Recommendation** | ✅ **Better for this project** | ❌ More complex |

## Troubleshooting

**"Invalid subscription key"**
- Make sure you copied the entire key
- Check there are no extra spaces
- Verify the resource is deployed

**"Out of call volume quota"**
- You've exceeded 1,000 calls this month
- Wait for next month or upgrade to S1 tier

**"Access denied"**
- Make sure you're using the Bing Search v7 API
- Not the old Bing Search v5 API

## Next Steps

After getting your key:

1. Add to `.env`:
   ```env
   NINJAS_API_KEY=your_ninjas_key
   BING_API_KEY=your_bing_key_here
   ```

2. Run the setup script:
   ```bash
   go run cmd/setup/main.go
   ```

3. Watch as it finds proper car images! 🚗
