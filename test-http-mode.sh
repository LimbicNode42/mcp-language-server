#!/bin/bash

# Quick test script for HTTP mode
# This creates a minimal test workspace and starts the server in HTTP mode

echo "Setting up test workspace..."
mkdir -p test-workspace
cd test-workspace

# Create a simple TypeScript file for testing
cat > test.ts << 'EOF'
function greet(name: string): string {
    return `Hello, ${name}!`;
}

const message = greet("World");
console.log(message);
EOF

# Create a simple package.json
cat > package.json << 'EOF'
{
    "name": "test-workspace",
    "version": "1.0.0",
    "scripts": {
        "build": "tsc"
    },
    "devDependencies": {
        "typescript": "^5.0.0"
    }
}
EOF

echo "Starting MCP Language Server in HTTP mode..."
echo "Server will be available at http://localhost:8080"
echo "Press Ctrl+C to stop"

cd ..
./mcp-language-server \
    --mode=http \
    --port=8080 \
    --workspace=./test-workspace \
    --lsp=typescript-language-server \
    -- --stdio
