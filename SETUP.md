# Quick Setup Guide

This guide helps you get started with the MCP Language Server using Docker.

## Prerequisites

- Docker and Docker Compose installed
- Your codebase/workspace ready

## Setup Steps

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/mcp-language-server.git
cd mcp-language-server
```

### 2. Create Environment Configuration

```bash
# Copy the example environment file
cp .env.example .env

# Edit the .env file to customize your setup
# At minimum, update the workspace paths to point to your code
```

### 3. Choose Your Deployment

#### Single Language Server

For TypeScript projects:
```bash
docker-compose up mcp-typescript
```

For Go projects:
```bash
docker-compose up mcp-go
```

For Python projects:
```bash
docker-compose up mcp-python
```

For Rust projects:
```bash
docker-compose up mcp-rust
```

For C/C++ projects:
```bash
docker-compose up mcp-clangd
```

#### Multiple Language Servers

```bash
# Start all language servers
docker-compose up

# Or specific ones
docker-compose up mcp-typescript mcp-go mcp-python
```

### 4. Configure Your MetaMCP Server

Update your MetaMCP server configuration to connect to the deployed services:

```json
{
  "mcpServers": {
    "typescript-language-server": {
      "transport": {
        "type": "http",
        "endpoint": "http://your-server:8080"
      }
    }
  }
}
```

## Default Ports

- TypeScript: 8080
- Go: 8081
- Python: 8082
- Rust: 8083
- C/C++ (clangd): 8084
- All-in-one: 8085

All ports are configurable via the `.env` file.

## Important Notes

- The `.env` file is not tracked in git - customize it for your setup
- Mount your actual codebase directories in the workspace paths
- Each service runs independently on its own port
- Health checks are included for monitoring
