# Freepik API Key Setup

Your Freepik API key has been configured. Here's how to set it up permanently:

## ✅ Quick Setup

### For Current Session (Already Done)
The API key is set for your current terminal session.

### For Permanent Setup (Local Development)

**Option 1: Add to your shell config file**

For Zsh (default on Mac):
```bash
echo 'export FREEPIK_API_KEY="FPSX762e644295c9849f0d26ff5df0b04c89"' >> ~/.zshrc
source ~/.zshrc
```

For Bash:
```bash
echo 'export FREEPIK_API_KEY="FPSX762e644295c9849f0d26ff5df0b04c89"' >> ~/.bashrc
source ~/.bashrc
```

**Option 2: Create a .env file (Recommended for Development) ✅ DONE**

A `.env` file has been created in your project root with the API key. The `start.sh` script will automatically load it when you run the server.

To use it:
```bash
./start.sh
```

Or if running manually:
```bash
export $(cat .env | xargs) && go run cmd/server/main.go
```

### For Production (Render.com)

1. Go to your Render.com dashboard: https://dashboard.render.com
2. Click on your `alice-suite-go` service
3. Go to the **Environment** tab
4. Click the **"Edit"** button next to "Environment Variables"
5. Add a new variable:
   - **KEY:** `FREEPIK_API_KEY`
   - **VALUE:** `FPSX762e644295c9849f0d26ff5df0b04c89`
6. Click **Save** and restart your service

## 🧪 Testing

To verify the API key is working:

1. Start your server:
   ```bash
   go run cmd/server/main.go
   ```

2. Log in as a reader
3. Open the AI Help assistant
4. Select some text and click "Visual Example"
5. You should see an image being generated (not just text description)

## 🔒 Security Note

⚠️ **Important:** Never commit your API key to Git. The `.gitignore` file already excludes `.env` files, so your key will be safe if you use the `.env` file method.

## 📝 What This Enables

With the Freepik API key configured, the "Visual Example" feature will:
- Generate actual images (PNG/JPEG) instead of just text descriptions
- Create pencil-style illustrations suitable for all ages
- Display images directly in the chat interface
- Allow users to click images to view full size

Without the key, visual examples will still work but only show text descriptions.
