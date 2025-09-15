# Quick test script for HTTP mode on Windows
# This creates a minimal test workspace and starts the server in HTTP mode

Write-Host "Setting up test workspace..."
New-Item -ItemType Directory -Force -Path "test-workspace" | Out-Null
Set-Location "test-workspace"

# Create a simple TypeScript file for testing
@"
function greet(name: string): string {
    return `Hello, `${name}!`;
}

const message = greet("World");
console.log(message);
"@ | Out-File -FilePath "test.ts" -Encoding UTF8

# Create a simple package.json
@"
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
"@ | Out-File -FilePath "package.json" -Encoding UTF8

Write-Host "Starting MCP Language Server in HTTP mode..."
Write-Host "Server will be available at http://localhost:8080"
Write-Host "Press Ctrl+C to stop"

Set-Location ".."
& ".\mcp-language-server.exe" `
    --mode=http `
    --port=8080 `
    --workspace=".\test-workspace" `
    --lsp=typescript-language-server `
    -- --stdio
