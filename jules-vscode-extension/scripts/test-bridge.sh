#!/bin/bash
# Test the Python bridge standalone

set -e

echo "🧪 Testing Python Bridge"
echo "======================="
echo ""

# Check if bridge exists
if [ ! -f "python/bridge.py" ]; then
    echo "❌ bridge.py not found. Run from extension root directory."
    exit 1
fi

# Check Python
if ! command -v python3 &> /dev/null; then
    echo "❌ Python 3 is not installed."
    exit 1
fi

# Check Jules SDK
if ! python3 -c "import jules_agent_sdk" 2>/dev/null; then
    echo "⚠️  Jules SDK not found. Installing..."
    pip install jules-agent-sdk
fi

echo "✅ Python environment OK"
echo ""
echo "Starting bridge in test mode..."
echo "Send JSON-RPC requests to test:"
echo ""
echo 'Example:'
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"api_key": "your-key"}}'
echo ""
echo "Press Ctrl+C to exit"
echo ""

python3 python/bridge.py
