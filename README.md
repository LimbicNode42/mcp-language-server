# MCP Language Server

[![Go Tests](https://github.com/isaacphi/mcp-language-server/actions/workflows/go.yml/badge.svg)](https://github.com/isaacphi/mcp-language-server/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/isaacphi/mcp-language-server)](https://goreportcard.com/report/github.com/isaacphi/mcp-language-server)
[![GoDoc](https://pkg.go.dev/badge/github.com/isaacphi/mcp-language-server)](https://pkg.go.dev/github.com/isaacphi/mcp-language-server)
[![Go Version](https://img.shields.io/github/go-mod/go-version/isaacphi/mcp-language-server)](https://github.com/isaacphi/mcp-language-server/blob/main/go.mod)

This is an [MCP](https://modelcontextprotocol.io/introduction) server that runs and exposes a [language server](https://microsoft.github.io/language-server-protocol/) to LLMs. Not a language server for MCP, whatever that would be.

## Demo

`mcp-language-server` helps MCP enabled clients navigate codebases more easily by giving them access semantic tools like get definition, references, rename, and diagnostics.

![Demo](demo.gif)

## Setup

1. **Install Go**: Follow instructions at <https://golang.org/doc/install>
2. **Install or update this server**: `go install github.com/isaacphi/mcp-language-server@latest`
3. **Install a language server**: _follow one of the guides below_
4. **Configure your MCP client**: _follow one of the guides below_

<details>
  <summary>Go (gopls)</summary>
  <div>
    <p><strong>Install gopls</strong>: <code>go install golang.org/x/tools/gopls@latest</code></p>
    <p><strong>Configure your MCP client</strong>: This will be different but similar for each client. For Claude Desktop, add the following to <code>~/Library/Application\ Support/Claude/claude_desktop_config.json</code></p>

<pre>
{
  "mcpServers": {
    "language-server": {
      "command": "mcp-language-server",
      "args": ["--workspace", "/Users/you/dev/yourproject/", "--lsp", "gopls"],
      "env": {
        "PATH": "/opt/homebrew/bin:/Users/you/go/bin",
        "GOPATH": "/users/you/go",
        "GOCACHE": "/users/you/Library/Caches/go-build",
        "GOMODCACHE": "/Users/you/go/pkg/mod"
      }
    }
  }
}
</pre>

<p><strong>Note</strong>: Not all clients will need these environment variables. For Claude Desktop you will need to update the environment variables above based on your machine and username:</p>
<ul>
  <li><code>PATH</code> needs to contain the path to <code>go</code> and to <code>gopls</code>. Get this with <code>echo $(which go):$(which gopls)</code></li>
  <li><code>GOPATH</code>, <code>GOCACHE</code>, and <code>GOMODCACHE</code> may be different on your machine. These are the defaults.</li>
</ul>

  </div>
</details>
<details>
  <summary>Rust (rust-analyzer)</summary>
  <div>
    <p><strong>Install rust-analyzer</strong>: <code>rustup component add rust-analyzer</code></p>
    <p><strong>Configure your MCP client</strong>: This will be different but similar for each client. For Claude Desktop, add the following to <code>~/Library/Application\ Support/Claude/claude_desktop_config.json</code></p>

<pre>
{
  "mcpServers": {
    "language-server": {
      "command": "mcp-language-server",
      "args": [
        "--workspace",
        "/Users/you/dev/yourproject/",
        "--lsp",
        "rust-analyzer"
      ]
    }
  }
}
</pre>
  </div>
</details>
<details>
  <summary>Python (pyright)</summary>
  <div>
    <p><strong>Install pyright</strong>: <code>npm install -g pyright</code></p>
    <p><strong>Configure your MCP client</strong>: This will be different but similar for each client. For Claude Desktop, add the following to <code>~/Library/Application\ Support/Claude/claude_desktop_config.json</code></p>

<pre>
{
  "mcpServers": {
    "language-server": {
      "command": "mcp-language-server",
      "args": [
        "--workspace",
        "/Users/you/dev/yourproject/",
        "--lsp",
        "pyright-langserver",
        "--",
        "--stdio"
      ]
    }
  }
}
</pre>
  </div>
</details>
<details>
  <summary>Typescript (typescript-language-server)</summary>
  <div>
    <p><strong>Install typescript-language-server</strong>: <code>npm install -g typescript typescript-language-server</code></p>
    <p><strong>Configure your MCP client</strong>: This will be different but similar for each client. For Claude Desktop, add the following to <code>~/Library/Application\ Support/Claude/claude_desktop_config.json</code></p>

<pre>
{
  "mcpServers": {
    "language-server": {
      "command": "mcp-language-server",
      "args": [
        "--workspace",
        "/Users/you/dev/yourproject/",
        "--lsp",
        "typescript-language-server",
        "--",
        "--stdio"
      ]
    }
  }
}
</pre>
  </div>
</details>
<details>
  <summary>C/C++ (clangd)</summary>
  <div>
    <p><strong>Install clangd</strong>: Download prebuilt binaries from the <a href="https://github.com/clangd/clangd/releases">official LLVM releases page</a> or install via your system's package manager (e.g., <code>apt install clangd</code>, <code>brew install clangd</code>).</p>
    <p><strong>Configure your MCP client</strong>: This will be different but similar for each client. For Claude Desktop, add the following to <code>~/Library/Application\\ Support/Claude/claude_desktop_config.json</code></p>

<pre>
{
  "mcpServers": {
    "language-server": {
      "command": "mcp-language-server",
      "args": [
        "--workspace",
        "/Users/you/dev/yourproject/",
        "--lsp",
        "/path/to/your/clangd_binary",
        "--",
        "--compile-commands-dir=/path/to/yourproject/build_or_compile_commands_dir"
      ]
    }
  }
}
</pre>
    <p><strong>Note</strong>:</p>
    <ul>
      <li>Replace <code>/path/to/your/clangd_binary</code> with the actual path to your clangd executable.</li>
      <li><code>--compile-commands-dir</code> should point to the directory containing your <code>compile_commands.json</code> file (e.g., <code>./build</code>, <code>./cmake-build-debug</code>).</li>
      <li>Ensure <code>compile_commands.json</code> is generated for your project for clangd to work effectively.</li>
    </ul>
  </div>
</details>
<details>
  <summary>Other</summary>
  <div>
    <p>I have only tested this repo with the servers above but it should be compatible with many more. Note:</p>
    <ul>
      <li>The language server must communicate over stdio.</li>
      <li>Any aruments after <code>--</code> are sent as arguments to the language server.</li>
      <li>Any env variables are passed on to the language server.</li>
    </ul>
  </div>
</details>

## HTTP Mode and Remote Deployment

In addition to the traditional stdio mode, this server now supports HTTP streaming mode for remote deployment. This is particularly useful for deploying language servers to cloud instances and accessing them from local development environments or MetaMCP servers.

### HTTP Mode Usage

Run the server in HTTP mode:

```bash
# TypeScript language server on port 8080
mcp-language-server --mode=http --port=8080 --workspace=/path/to/project --lsp=typescript-language-server -- --stdio

# Go language server on port 8081  
mcp-language-server --mode=http --port=8081 --workspace=/path/to/project --lsp=gopls

# Python language server on port 8082
mcp-language-server --mode=http --port=8082 --workspace=/path/to/project --lsp=pyright-langserver -- --stdio
```

### Docker Deployment

The easiest way to deploy remotely is using Docker. Several deployment options are provided:

#### Quick Start with TypeScript

```bash
# Clone and build
git clone https://github.com/your-username/mcp-language-server.git
cd mcp-language-server

# Copy and customize environment variables
cp .env.example .env
# Edit .env to set your workspace paths and ports

# Start TypeScript language server
docker-compose up mcp-typescript
```

#### Multi-Language Deployment

Deploy multiple language servers simultaneously:

```bash
# Start all language servers on different ports
docker-compose up mcp-typescript mcp-go mcp-python mcp-rust mcp-clangd
```

#### Individual Language Servers

Deploy specific language servers:

```bash
# TypeScript only
docker-compose up mcp-typescript

# Go only  
docker-compose up mcp-go

# Python only
docker-compose up mcp-python

# Rust only
docker-compose up mcp-rust

# C/C++ (clangd) only
docker-compose up mcp-clangd
```

#### Environment Configuration

All aspects of the deployment are configurable via `.env` file. Key environment variables:

```bash
# Language server ports
MCP_TYPESCRIPT_PORT=8080
MCP_GO_PORT=8081
MCP_PYTHON_PORT=8082
MCP_RUST_PORT=8083
MCP_CLANGD_PORT=8084

# Workspace paths (mount your code here)
WORKSPACE_PATH=./workspace
GO_WORKSPACE_PATH=./go-workspace
PYTHON_WORKSPACE_PATH=./python-workspace
RUST_WORKSPACE_PATH=./rust-workspace
CLANGD_WORKSPACE_PATH=./clangd-workspace

# Service-specific configuration (all customizable)
TYPESCRIPT_MCP_MODE=http
GO_MCP_MODE=http
PYTHON_MCP_MODE=http
# ... and many more - see .env.example for full list
```

**Complete configurability**: Every environment variable used in the docker-compose.yml can be customized via the .env file, including service modes, internal ports, workspace directories, and language server types.

### Connecting from MetaMCP Server

Configure your MetaMCP server to connect to the deployed language servers:

```json
{
  "mcpServers": {
    "typescript-language-server": {
      "transport": {
        "type": "http",
        "endpoint": "http://your-server:8080"
      }
    },
    "go-language-server": {
      "transport": {
        "type": "http", 
        "endpoint": "http://your-server:8081"
      }
    }
  }
}
```

### Reverse Proxy Configuration

If you want to use a reverse proxy (like Traefik), configure it to route to the individual service endpoints:

- TypeScript: `http://mcp-typescript:8080` 
- Go: `http://mcp-go:8080`
- Python: `http://mcp-python:8080`
- Rust: `http://mcp-rust:8080`
- C/C++: `http://mcp-clangd:8080`

### Production Considerations

- **Workspace Security**: Mount only the specific directories you need
- **Network Security**: Use proper firewall rules and consider VPN/private networks
- **Resource Limits**: Configure Docker memory/CPU limits for large codebases
- **Persistent Storage**: Use Docker volumes for caches (Go build cache, Cargo cache, etc.)
- **Reverse Proxy**: Use your existing reverse proxy (like Traefik) to route to individual services
- **Health Monitoring**: Health check endpoints are included in the Docker configurations

## Tools

- `definition`: Retrieves the complete source code definition of any symbol (function, type, constant, etc.) from your codebase.
- `references`: Locates all usages and references of a symbol throughout the codebase.
- `diagnostics`: Provides diagnostic information for a specific file, including warnings and errors.
- `hover`: Display documentation, type hints, or other hover information for a given location.
- `rename_symbol`: Rename a symbol across a project.
- `edit_file`: Allows making multiple text edits to a file based on line numbers. Provides a more reliable and context-economical way to edit files compared to search and replace based edit tools.

## About

This codebase makes use of edited code from [gopls](https://go.googlesource.com/tools/+/refs/heads/master/gopls/internal/protocol) to handle LSP communication. See ATTRIBUTION for details. Everything here is covered by a permissive BSD style license.

The [official MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) is used for MCP communication, providing both stdio and HTTP streaming transport support.

This is beta software. Please let me know by creating an issue if you run into any problems or have suggestions of any kind.

## Contributing

Please keep PRs small and open Issues first for anything substantial. AI slop O.K. as long as it is tested, passes checks, and doesn't smell too bad.

### Setup

Clone the repo:

```bash
git clone https://github.com/isaacphi/mcp-language-server.git
cd mcp-language-server
```

A [justfile](https://just.systems/man/en/) is included for convenience:

```bash
just -l
Available recipes:
    build    # Build
    check    # Run code audit checks
    fmt      # Format code
    generate # Generate LSP types and methods
    help     # Help
    install  # Install locally
    snapshot # Update snapshot tests
    test     # Run tests
```

Configure your Claude Desktop (or similar) to use the local binary:

```json
{
  "mcpServers": {
    "language-server": {
      "command": "/full/path/to/your/clone/mcp-language-server/mcp-language-server",
      "args": [
        "--workspace",
        "/path/to/workspace",
        "--lsp",
        "language-server-executable"
      ],
      "env": {
        "LOG_LEVEL": "DEBUG"
      }
    }
  }
}
```

Rebuild after making changes.

### Logging

Setting the `LOG_LEVEL` environment variable to DEBUG enables verbose logging to stderr for all components including messages to and from the language server and the language server's logs.

### LSP interaction

- `internal/lsp/methods.go` contains generated code to make calls to the connected language server.
- `internal/protocol/tsprotocol.go` contains generated code for LSP types. I borrowed this from `gopls`'s source code. Thank you for your service.
- LSP allows language servers to return different types for the same methods. Go doesn't like this so there are some ugly workarounds in `internal/protocol/interfaces.go`.

### Local Development and Snapshot Tests

There is a snapshot test suite that makes it a lot easier to try out changes to tools. These run actual language servers on mock workspaces and capture output and logs.

You will need the language servers installed locally to run them. There are tests for go, rust, python, and typescript.

```
integrationtests/
├── tests/        # Tests are in this folder
├── snapshots/    # Snapshots of tool outputs
├── test-output/  # Gitignored folder showing the final state of each workspace and logs after each test run
└── workspaces/   # Mock workspaces that the tools run on
```

To update snapshots, run `UPDATE_SNAPSHOTS=true go test ./integrationtests/...`
