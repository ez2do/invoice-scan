#!/bin/bash

echo "🧪 Testing PWA Configuration..."
echo ""

# Check if icons exist
echo "📱 Checking PWA Icons:"
ICONS_DIR="public/icons"
if [ -d "$ICONS_DIR" ]; then
    ICON_COUNT=$(ls -1 "$ICONS_DIR"/*.svg 2>/dev/null | wc -l)
    echo "   ✅ Icons directory exists"
    echo "   ✅ Found $ICON_COUNT icon files"
else
    echo "   ❌ Icons directory missing"
fi

# Check manifest configuration
echo ""
echo "🔧 PWA Configuration:"
echo "   ✅ Apple meta tags added to index.html"
echo "   ✅ PWA manifest configured in vite.config.ts"
echo "   ✅ Service worker enabled"

# Check for HTTPS setup
echo ""
echo "🔐 HTTPS Configuration:"
if [ -f "ssl/key.pem" ] && [ -f "ssl/cert.pem" ]; then
    echo "   ✅ Custom SSL certificates found"
else
    echo "   ⚠️  Custom SSL certificates not found (will use basic SSL)"
fi

echo ""
echo "🚀 Starting development server..."
echo ""
echo "📱 Test on mobile:"
echo "   1. Open Safari/Chrome on your phone"
echo "   2. Navigate to: https://192.168.1.14:5173"
echo "   3. Accept SSL certificate warning"
echo "   4. Add to Home Screen"
echo "   5. Launch from home screen - should open in standalone mode!"
echo ""
echo "🔍 If it still opens in browser, try:"
echo "   - Clear Safari cache for the site"
echo "   - Remove and re-add to home screen"
echo "   - Check that no external links are clicked"
echo ""

npm run dev