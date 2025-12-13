#!/bin/bash

# Safari Proxy Configuration Script
# This configures Safari to bypass proxy for localhost
# Run this once, or add to your .zshrc to run automatically

echo "=== Configuring Safari Proxy Bypass ==="
echo ""

# Detect network interface (try multiple methods)
INTERFACE=$(networksetup -listallnetworkservices 2>/dev/null | grep -E "Wi-Fi|Ethernet" | head -1)

# If that fails, try to get active interface
if [ -z "$INTERFACE" ]; then
    INTERFACE=$(route get default 2>/dev/null | grep interface | awk '{print $2}' | xargs networksetup -listnetworkserviceorder 2>/dev/null | grep -A1 "$(route get default 2>/dev/null | grep interface | awk '{print $2}')" | head -1 | sed 's/.*: //')
fi

# Last resort: try common interface names
if [ -z "$INTERFACE" ]; then
    for iface in "Wi-Fi" "Ethernet" "en0" "en1"; do
        if networksetup -getproxybypassdomains "$iface" >/dev/null 2>&1; then
            INTERFACE="$iface"
            break
        fi
    done
fi

if [ -z "$INTERFACE" ]; then
    echo "❌ Could not detect network interface"
    echo "Available network services:"
    networksetup -listallnetworkservices 2>/dev/null || echo "Could not list services"
    echo ""
    echo "Please run manually:"
    echo "  sudo networksetup -setproxybypassdomains 'Wi-Fi' '127.0.0.1' 'localhost' '*.local'"
    exit 1
fi

echo "Network interface: $INTERFACE"
echo ""

# Get current bypass domains
CURRENT_BYPASS=$(networksetup -getproxybypassdomains "$INTERFACE" 2>/dev/null | tr '\n' ' ' || echo "")

# Check if localhost is already configured
if echo "$CURRENT_BYPASS" | grep -q "127.0.0.1\|localhost"; then
    echo "✅ localhost is already in proxy bypass list"
    echo "Current bypass domains:"
    networksetup -getproxybypassdomains "$INTERFACE" | sed 's/^/   /'
else
    echo "⚠️  localhost is NOT in proxy bypass list"
    echo ""
    echo "Adding localhost to proxy bypass..."
    echo "(This requires your admin password)"
    echo ""
    
    sudo networksetup -setproxybypassdomains "$INTERFACE" \
        "127.0.0.1" "localhost" "*.local" "169.254/16"
    
    if [ $? -eq 0 ]; then
        echo ""
        echo "✅ Successfully configured Safari proxy bypass!"
        echo ""
        echo "Updated bypass domains:"
        networksetup -getproxybypassdomains "$INTERFACE" | sed 's/^/   /'
        echo ""
        echo "Safari will now:"
        echo "  ✅ Connect directly to localhost/127.0.0.1 (no proxy)"
        echo "  ✅ Use proxy for internet sites (with your credentials)"
    else
        echo ""
        echo "❌ Failed to update proxy bypass"
        echo ""
        echo "Manual configuration:"
        echo "  1. Safari → Preferences → Advanced"
        echo "  2. Click 'Change Settings...' next to 'Proxies'"
        echo "  3. In 'Bypass proxy settings', add: 127.0.0.1, localhost"
        exit 1
    fi
fi

echo ""
echo "=== Configuration Complete ==="

