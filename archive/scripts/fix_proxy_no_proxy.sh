#!/bin/bash

# Fix Proxy Configuration to Bypass Localhost
# This script helps configure NO_PROXY to include localhost

echo "=== Fixing Proxy NO_PROXY Configuration ==="
echo ""

# Check current proxy settings
echo "Current proxy settings:"
env | grep -i proxy | grep -v PASS
echo ""

# Check if no_proxy is set
if [ -z "$no_proxy" ] && [ -z "$NO_PROXY" ]; then
    echo "⚠️  NO_PROXY is not set - localhost will go through proxy!"
    echo ""
    echo "To fix this, add to your ~/.secure/proxy_config.sh:"
    echo ""
    echo "  export no_proxy=\"127.0.0.1,localhost,*.local,0.0.0.0\""
    echo "  export NO_PROXY=\"127.0.0.1,localhost,*.local,0.0.0.0\""
    echo ""
else
    echo "Current NO_PROXY: $no_proxy $NO_PROXY"
    if echo "$no_proxy $NO_PROXY" | grep -q "127.0.0.1\|localhost"; then
        echo "✅ localhost is already in NO_PROXY"
    else
        echo "⚠️  localhost is NOT in NO_PROXY - needs to be added"
    fi
fi

echo ""
echo "=== Checking ~/.secure/proxy_config.sh ==="
if [ -f ~/.secure/proxy_config.sh ]; then
    echo "File exists. Current content:"
    echo "---"
    cat ~/.secure/proxy_config.sh
    echo "---"
    echo ""
    
    if grep -q "no_proxy\|NO_PROXY" ~/.secure/proxy_config.sh; then
        echo "NO_PROXY is configured in the file"
    else
        echo "⚠️  NO_PROXY is NOT configured - adding it..."
        echo ""
        read -p "Do you want to add NO_PROXY to ~/.secure/proxy_config.sh? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "" >> ~/.secure/proxy_config.sh
            echo "# Bypass proxy for localhost and local domains" >> ~/.secure/proxy_config.sh
            echo "export no_proxy=\"127.0.0.1,localhost,*.local,0.0.0.0\"" >> ~/.secure/proxy_config.sh
            echo "export NO_PROXY=\"127.0.0.1,localhost,*.local,0.0.0.0\"" >> ~/.secure/proxy_config.sh
            echo ""
            echo "✅ Added NO_PROXY to ~/.secure/proxy_config.sh"
            echo ""
            echo "Now reload proxy settings:"
            echo "  proxy_off"
            echo "  proxy_on"
        fi
    fi
else
    echo "⚠️  ~/.secure/proxy_config.sh not found"
    echo ""
    echo "You can create it with:"
    echo "  mkdir -p ~/.secure"
    echo "  cat > ~/.secure/proxy_config.sh << 'EOF'"
    echo "export http_proxy=\"http://nepepittau:app730@192.168.77.1:8000\""
    echo "export https_proxy=\"\$http_proxy\""
    echo "export HTTP_PROXY=\"\$http_proxy\""
    echo "export HTTPS_PROXY=\"\$http_proxy\""
    echo "export no_proxy=\"127.0.0.1,localhost,*.local,0.0.0.0\""
    echo "export NO_PROXY=\"127.0.0.1,localhost,*.local,0.0.0.0\""
    echo "EOF"
fi

echo ""
echo "=== Important Notes ==="
echo "1. NO_PROXY in shell only affects command-line tools (curl, git, etc.)"
echo "2. Safari uses macOS system proxy settings (not shell variables)"
echo "3. To fix Safari, you still need to:"
echo "   - Run: sudo networksetup -setproxybypassdomains 'Wi-Fi' '127.0.0.1' 'localhost' '*.local'"
echo "   - OR configure Safari → Preferences → Advanced → Proxies → Bypass"
echo ""

