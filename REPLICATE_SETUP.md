# Replicate Image Generation Setup

Replicate has been integrated as an image generation provider option. It offers access to Stable Diffusion models and other AI image generation models.

## Setup

### 1. Get Your Replicate API Token

1. Go to [Replicate](https://replicate.com)
2. Sign up or log in
3. Go to your account settings: https://replicate.com/account/api-tokens
4. Copy your API token

### 2. Set Environment Variable

**Local Development (.env file):**
```bash
REPLICATE_API_TOKEN=your-replicate-api-token-here
IMAGE_PROVIDER=replicate
```

**Or export in terminal:**
```bash
export REPLICATE_API_TOKEN="your-replicate-api-token-here"
export IMAGE_PROVIDER="replicate"
```

**Render.com:**
1. Go to Environment tab
2. Add: `REPLICATE_API_TOKEN` = `your-replicate-api-token-here`
3. Add: `IMAGE_PROVIDER` = `replicate` (optional, defaults to deepai)

### 3. Restart Server

After setting the environment variables, restart your server:
```bash
./start.sh
```

## Features

- **Stable Diffusion models**: Uses Stable Diffusion 3 by default (high quality)
- **Automatic pencil sketch styling**: Prompts are automatically enhanced with "minimalist pencil sketch" keywords
- **Asynchronous generation**: Returns prediction ID, then polls for completion
- **Educational illustrations**: Perfect for pencil-style educational content

## Default Model

The default model is `stability-ai/stable-diffusion-3`. You can specify a different model in the request if needed.

## Switching Providers

Set `IMAGE_PROVIDER` environment variable:
- `replicate` - Uses Replicate (Stable Diffusion models)
- `deepai` - Uses DeepAI (default, better free tier)
- `freepik` - Uses Freepik (requires FREEPIK_API_KEY)

## Pricing

Replicate uses a pay-as-you-go pricing model. Check [Replicate Pricing](https://replicate.com/pricing) for current rates.
