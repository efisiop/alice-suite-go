#!/bin/bash

# Safari Proxy Configuration Script
# This script helps configure Safari to bypass proxy for localhost
# while keeping proxy enabled for internet access

echo "=== Safari Proxy Configuration Helper ==="
echo ""
echo "Your current proxy settings:"
echo "  Proxy: $http_proxy"
echo ""

echo "⚠️  IMPORTANT:"
echo "   localhost (127.0.0.1) should NOT go through the proxy!"
echo "   The proxy is for internet access, not local development."
echo ""
echo "To configure Safari properly:"
echo ""
echo "1. Open Safari"
echo "2. Go to: Safari → Preferences (or Settings)"
echo "3. Click the 'Advanced' tab"
echo "4. Click 'Change Settings...' next to 'Proxies'"
echo "5. In the 'Bypass proxy settings for these Hosts & Domains' field, add:"
echo ""
echo "   127.0.0.1, localhost, *.local, 0.0.0.0"
echo ""
echo "6. Make sure your proxy settings are still configured for internet access"
echo "7. Click OK and close preferences"
echo ""
echo "This way:"
echo "  ✅ localhost/127.0.0.1 → Direct connection (no proxy)"
echo "  ✅ Internet sites → Through proxy (with your credentials)"
echo ""
echo "Alternative: Use networksetup command (requires admin):"
echo ""
echo "  sudo networksetup -setwebproxystate 'Wi-Fi' on"
echo "  sudo networksetup -setwebproxy 'Wi-Fi' 192.168.77.1 8000"
echo "  sudo networksetup -setproxybypassdomains 'Wi-Fi' '127.0.0.1' 'localhost' '*.local'"
echo ""
echo "To check current proxy bypass settings:"
echo "  networksetup -getproxybypassdomains 'Wi-Fi'"
echo ""

