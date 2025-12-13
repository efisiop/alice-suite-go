#!/bin/bash

# Update Proxy Config to Include NO_PROXY for Localhost
# This updates ~/.secure/proxy_config.sh to uncomment and configure NO_PROXY

PROXY_FILE="$HOME/.secure/proxy_config.sh"

if [ ! -f "$PROXY_FILE" ]; then
    echo "❌ Proxy config file not found: $PROXY_FILE"
    exit 1
fi

echo "=== Updating Proxy Configuration ==="
echo ""

# Backup the original file
cp "$PROXY_FILE" "${PROXY_FILE}.backup.$(date +%Y%m%d_%H%M%S)"
echo "✅ Created backup: ${PROXY_FILE}.backup.*"
echo ""

# Check if NO_PROXY is already uncommented
if grep -q "^export no_proxy=" "$PROXY_FILE"; then
    echo "✅ NO_PROXY is already configured"
    echo "Current NO_PROXY settings:"
    grep "^export.*no_proxy\|^export.*NO_PROXY" "$PROXY_FILE"
else
    echo "Updating NO_PROXY configuration..."
    
    # Uncomment and update NO_PROXY lines
    sed -i '' 's|^# export no_proxy=|export no_proxy=|' "$PROXY_FILE"
    sed -i '' 's|^# export NO_PROXY=|export NO_PROXY=|' "$PROXY_FILE"
    
    # Update the values to include localhost
    sed -i '' 's|export no_proxy=".*"|export no_proxy="127.0.0.1,localhost,*.local,0.0.0.0"|' "$PROXY_FILE"
    sed -i '' 's|export NO_PROXY=.*|export NO_PROXY="127.0.0.1,localhost,*.local,0.0.0.0"|' "$PROXY_FILE"
    
    echo "✅ Updated NO_PROXY in $PROXY_FILE"
fi

echo ""
echo "Updated file content:"
echo "---"
grep -A 2 "NO_PROXY\|no_proxy" "$PROXY_FILE" | head -5
echo "---"
echo ""
echo "To apply changes:"
echo "  1. Reload proxy: proxy_off && proxy_on"
echo "  2. Verify: echo \$no_proxy"
echo ""
echo "⚠️  Remember: Safari needs separate configuration!"
echo "  Run: sudo networksetup -setproxybypassdomains 'Wi-Fi' '127.0.0.1' 'localhost' '*.local'"

