#!/bin/bash

echo "🔐 Setting up SSL for mobile development..."

# Check if mkcert is installed
if command -v mkcert &> /dev/null; then
    echo "✅ mkcert found, creating certificates..."
    
    # Create ssl directory
    mkdir -p ssl
    
    # Generate certificates for your local IP
    mkcert -key-file ssl/key.pem -cert-file ssl/cert.pem 192.168.1.14 localhost 127.0.0.1
    
    echo "✅ SSL certificates created in ssl/ directory"
    echo "📱 You can now use https://192.168.1.14:5173 on mobile"
    
else
    echo "❌ mkcert not found. Installing..."
    
    # Detect OS and install mkcert
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            brew install mkcert
            mkcert -install
        else
            echo "Please install Homebrew first: https://brew.sh"
            exit 1
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        echo "Please install mkcert manually:"
        echo "curl -JLO 'https://dl.filippo.io/mkcert/latest?for=linux/amd64'"
        echo "chmod +x mkcert-v*-linux-amd64"
        echo "sudo cp mkcert-v*-linux-amd64 /usr/local/bin/mkcert"
        echo "mkcert -install"
        exit 1
    else
        echo "Unsupported OS. Please install mkcert manually: https://github.com/FiloSottile/mkcert"
        exit 1
    fi
    
    # Try again after installation
    if command -v mkcert &> /dev/null; then
        mkdir -p ssl
        mkcert -key-file ssl/key.pem -cert-file ssl/cert.pem 192.168.1.14 localhost 127.0.0.1
        echo "✅ SSL certificates created"
    fi
fi

echo ""
echo "🚀 Next steps:"
echo "1. Run: npm run dev"
echo "2. Open https://192.168.1.14:5173 on mobile"
echo "3. Accept certificate when prompted"