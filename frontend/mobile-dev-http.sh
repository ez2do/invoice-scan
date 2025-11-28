#!/bin/bash

echo "📱 Starting PWA for mobile development (HTTP workaround)..."
echo ""
echo "⚠️  CAMERA LIMITATION: Chrome requires HTTPS for camera access"
echo "   This HTTP version is for UI testing only"
echo ""
echo "📱 Access on mobile:"
echo "   http://192.168.1.14:5173"
echo ""
echo "🔧 For camera functionality, use one of these options:"
echo "   1. Run ./setup-ssl.sh to create proper HTTPS"
echo "   2. Use Chrome's --unsafely-treat-insecure-origin-as-secure flag"
echo "   3. Test on desktop Chrome with http://localhost:5173"
echo ""
echo "Starting development server..."

# Temporarily remove basicSsl plugin for HTTP-only mode
npm run dev