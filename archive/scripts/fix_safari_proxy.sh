#!/bin/bash

# Fix Safari Proxy to Bypass Localhost
# This adds localhost to the proxy bypass list while keeping proxy for internet

echo "=== Fixing Safari Proxy Configuration ==="
echo ""

# Detect network interface
INTERFACE=$(networksetup -listallnetworkservices | grep -E "Wi-Fi|Ethernet" | head -1)
if [ -z "$INTERFACE" ]; then
    echo "❌ Could not detect network interface"
    echo "Please run manually:"
    echo "  networksetup -getproxybypassdomains 'Wi-Fi'"
    exit 1
fi

echo "Detected network interface: $INTERFACE"
echo ""

# Get current bypass domains
CURRENT_BYPASS=$(networksetup -getproxybypassdomains "$INTERFACE" 2>/dev/null | tr '\n' ' ')

echo "Current bypass domains: $CURRENT_BYPASS"
echo ""

# Check if localhost is already in bypass list
if echo "$CURRENT_BYPASS" | grep -q "127.0.0.1\|localhost"; then
    echo "✅ localhost is already in bypass list"
else
    echo "⚠️  localhost is NOT in bypass list - adding it..."
    echo ""
    echo "To add localhost to proxy bypass (requires admin password):"
    echo ""
    echo "  sudo networksetup -setproxybypassdomains '$INTERFACE' \\"
    echo "    '127.0.0.1' 'localhost' '*.local' '169.254/16'"
    echo ""
    echo "Or add manually in Safari:"
    echo "  Safari → Preferences → Advanced → Proxies → Bypass proxy settings"
    echo "  Add: 127.0.0.1, localhost"
    echo ""
    
    # Try to add it (will prompt for password)
    read -p "Do you want to add localhost to bypass list now? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        sudo networksetup -setproxybypassdomains "$INTERFACE" \
            "127.0.0.1" "localhost" "*.local" "169.254/16"
        
        if [ $? -eq 0 ]; then
            echo ""
            echo "✅ Successfully added localhost to proxy bypass!"
            echo ""
            echo "New bypass domains:"
            networksetup -getproxybypassdomains "$INTERFACE"
            echo ""
            echo "Now Safari should work with http://127.0.0.1:8080"
        else
            echo ""
            echo "❌ Failed to update proxy bypass settings"
        fi
    fi
fi

echo ""
echo "=== Test Connection ==="
echo "Testing if server is accessible..."
if curl --noproxy "*" -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8080/health | grep -q "200"; then
    echo "✅ Server is running and accessible"
else
    echo "❌ Server is not accessible - make sure it's running"
fi

