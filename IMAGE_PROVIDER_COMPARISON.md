# Image Generation Provider Comparison

## Current Setup

**Active Provider:** DeepAI (free tier)
**Configuration:** Set in `.env` file via `IMAGE_PROVIDER` variable

## Provider Options

### 1. DeepAI (Current - Recommended for Free Tier)
- **Free Tier:** ~100 images/month
- **API Key:** `DEEPAI_API_KEY`
- **Setup:** Get key from https://deepai.org/api-key
- **Pros:** 
  - Generous free tier
  - Synchronous (fast response)
  - Simple setup
- **Cons:**
  - Limited to 100 images/month on free tier
- **Best for:** Development, testing, low-volume usage

### 2. Replicate (Requires Credits)
- **Pricing:** Pay-as-you-go (requires purchasing credits)
- **API Token:** `REPLICATE_API_TOKEN`
- **Setup:** Get token from https://replicate.com/account/api-tokens
- **Pros:**
  - High quality Stable Diffusion models
  - Flexible model selection
  - Good for production
- **Cons:**
  - Requires purchasing credits upfront
  - No free tier
- **Best for:** Production with budget, when quality is critical

### 3. Freepik (Issues Encountered)
- **API Key:** `FREEPIK_API_KEY`
- **Pros:**
  - Dedicated image generation service
- **Cons:**
  - TLS certificate issues encountered
  - Timeout issues
  - API reliability concerns
- **Status:** Not recommended until issues are resolved

## Switching Providers

### To Use DeepAI (Current - Free):
```bash
export IMAGE_PROVIDER=deepai
# Or in .env file:
IMAGE_PROVIDER=deepai
DEEPAI_API_KEY=your-key-here
```

### To Use Replicate (Requires Credits):
```bash
export IMAGE_PROVIDER=replicate
# Or in .env file:
IMAGE_PROVIDER=replicate
REPLICATE_API_TOKEN=your-token-here
```

**Note:** Before using Replicate, make sure to:
1. Go to https://replicate.com/account/billing
2. Purchase credits
3. Wait a few minutes after purchase
4. Then restart your server

## Recommendation

For now, **DeepAI is the best choice** because:
- ✅ Free tier available
- ✅ No credit card required
- ✅ Works reliably
- ✅ Good for educational illustrations

You can switch to Replicate later when you're ready to purchase credits for higher quality/volume.
