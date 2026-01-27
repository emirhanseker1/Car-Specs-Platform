# How to Get SerpApi Key (Google Images)

## Overview
SerpApi provides easy access to Google Search results, including Google Images. It's the simplest solution - no Azure subscription, no complex setup, just sign up and get an API key.

## Why SerpApi?
- ✅ **100 free searches/month** (no credit card required)
- ✅ **Instant setup** (< 2 minutes)
- ✅ **Direct Google Images access**
- ✅ **No subscription or billing setup needed**
- ✅ **Simple JSON API**

## Step-by-Step Setup

### 1. Sign Up for SerpApi

1. **Go to SerpApi**
   - Visit: https://serpapi.com/

2. **Create Free Account**
   - Click "Register" or "Sign Up"
   - Enter your email and password
   - Or sign up with Google/GitHub

3. **Verify Email**
   - Check your email for verification link
   - Click to verify your account

### 2. Get Your API Key

1. **Go to Dashboard**
   - After login, you'll see your dashboard
   - Or visit: https://serpapi.com/manage-api-key

2. **Copy Your API Key**
   - You'll see "Your Private API Key"
   - It looks like: `a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0`
   - Click "Copy" or select and copy manually

3. **Add to .env File**
   ```env
   SERPAPI_KEY=your_actual_key_here
   ```

### 3. Free Tier Details

**What you get FREE:**
- 100 searches per month
- Access to Google Images
- Access to Google Search
- JSON API responses
- No credit card required

**What happens after 100 searches:**
- Your requests will fail
- You can upgrade to paid plan ($50/month for 5,000 searches)
- Or wait for next month (resets monthly)

## Test Your API Key

Test it with curl:

```bash
curl "https://serpapi.com/search.json?q=bmw+320i+2023+car&tbm=isch&api_key=YOUR_KEY_HERE"
```

You should get a JSON response with image results.

## Usage in Our Script

The setup script uses SerpApi like this:

```go
// Query format
query := "BMW 3 Series 320i 2023 exterior car studio white background"

// API call
params := url.Values{}
params.Add("q", query)
params.Add("tbm", "isch")  // Target: Images
params.Add("api_key", os.Getenv("SERPAPI_KEY"))

// Response
{
  "images_results": [
    {
      "original": "https://example.com/car-image.jpg",
      "thumbnail": "https://...",
      "title": "BMW 320i 2023"
    }
  ]
}
```

## Comparison: SerpApi vs Others

| Feature | SerpApi | Bing API | Google Custom Search |
|---------|---------|----------|---------------------|
| Free Tier | 100/month | 1,000/month | 100/day |
| Setup Time | **2 minutes** | 15 minutes | 20 minutes |
| Credit Card | **Not required** | Not required | Required |
| Complexity | **Very simple** | Medium | Complex |
| Image Quality | **Excellent** | Excellent | Excellent |
| **Best For** | **Quick setup** | High volume | Custom filtering |

## Troubleshooting

**"Invalid API key"**
- Make sure you copied the entire key
- Check there are no extra spaces
- Verify you're logged into SerpApi

**"You have reached your monthly limit"**
- You've used all 100 free searches
- Wait for next month or upgrade plan

**"No results found"**
- The search query might be too specific
- Try a simpler query
- Check if the car model exists

## Pricing (Optional)

If you need more than 100 searches/month:

- **Starter**: $50/month - 5,000 searches
- **Developer**: $100/month - 15,000 searches
- **Business**: Custom pricing

For this project, **100 free searches is plenty** for initial setup.

## Next Steps

After getting your key:

1. Add to `.env`:
   ```env
   NINJAS_API_KEY=your_ninjas_key
   SERPAPI_KEY=your_serpapi_key_here
   ```

2. Run the setup script:
   ```bash
   go run cmd/setup/main.go
   ```

3. Watch as it finds proper car images from Google! 🚗

## Alternative: Run Without Images

If you don't want to sign up for SerpApi, the script works fine without it:

```env
NINJAS_API_KEY=your_ninjas_key
# SERPAPI_KEY=  (leave empty or comment out)
```

The script will:
- ✅ Still fetch all vehicle data
- ✅ Still create brands, models, trims
- ⚠️ Just skip image search
- You can add images later manually or with a different service
