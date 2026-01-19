# DeepAI Image Generation Setup

DeepAI has been integrated as the default image generation provider. It offers:
- **Free tier**: ~100 images/month
- **Simple API**: Synchronous (returns image immediately)
- **Pencil sketch support**: Great for educational illustrations

## Setup

### 1. Get Your DeepAI API Key

1. Go to [DeepAI API](https://deepai.org/api-key)
2. Sign up or log in (free account)
3. Copy your API key

### 2. Set Environment Variable

**Local Development (.env file):**
```bash
DEEPAI_API_KEY=your-deepai-api-key-here
```

**Or export in terminal:**
```bash
export DEEPAI_API_KEY="your-deepai-api-key-here"
```

**Render.com:**
1. Go to Environment tab
2. Add: `DEEPAI_API_KEY` = `your-deepai-api-key-here`

### 3. Configure Provider (Optional)

By default, DeepAI is used. To switch back to Freepik:
```bash
export IMAGE_PROVIDER="freepik"
```

## Features

- **Automatic pencil sketch styling**: Prompts are automatically enhanced with "minimalist pencil sketch" keywords
- **Synchronous generation**: Images return immediately (no polling needed)
- **Free tier friendly**: 100 images/month for free accounts

## Switching Providers

Set `IMAGE_PROVIDER` environment variable:
- `deepai` - Default, better free tier
- `freepik` - Requires FREEPIK_API_KEY
