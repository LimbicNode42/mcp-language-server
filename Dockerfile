# Multi-stage Dockerfile for MCP Language Server
# Supports TypeScript, Go, Python, Rust, and C/C++ language servers

FROM node:20-alpine AS node-base
RUN npm install -g typescript typescript-language-server pyright

FROM golang:1.22-alpine AS go-base
RUN apk add --no-cache git
RUN go install golang.org/x/tools/gopls@latest

FROM rust:1.75-alpine AS rust-base
RUN apk add --no-cache git
RUN rustup component add rust-analyzer

FROM alpine:3.18 AS clangd-base
RUN apk add --no-cache clang clang-extra-tools

# Build stage for mcp-language-server
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mcp-language-server .

# Base runtime image
FROM alpine:3.18 AS runtime-base
RUN apk add --no-cache \
    ca-certificates \
    git \
    python3 \
    py3-pip \
    nodejs \
    npm

# Final stage with all language servers
FROM runtime-base AS all-languages

# Copy Node.js language servers
COPY --from=node-base /usr/local/lib/node_modules /usr/local/lib/node_modules
COPY --from=node-base /usr/local/bin/tsserver /usr/local/bin/
COPY --from=node-base /usr/local/bin/typescript-language-server /usr/local/bin/
COPY --from=node-base /usr/local/bin/pyright-langserver /usr/local/bin/

# Copy Go language server
COPY --from=go-base /go/bin/gopls /usr/local/bin/

# Copy Rust language server
COPY --from=rust-base /usr/local/rustup/toolchains/*/bin/rust-analyzer /usr/local/bin/

# Copy clangd
COPY --from=clangd-base /usr/bin/clangd /usr/local/bin/

# Copy the built mcp-language-server
COPY --from=builder /app/mcp-language-server /usr/local/bin/

# Create workspace directory
RUN mkdir -p /workspace
WORKDIR /workspace

# Environment variables for configuration
ENV LANGUAGE_SERVER_TYPE=typescript
ENV MCP_PORT=8080
ENV MCP_MODE=http
ENV WORKSPACE_DIR=/workspace

# Expose the port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD if [ "$MCP_MODE" = "http" ]; then \
        wget --no-verbose --tries=1 --spider http://localhost:${MCP_PORT}/ || exit 1; \
    else \
        echo "Stdio mode - skipping health check"; \
    fi

# Default command - can be overridden via environment variables
CMD case "$LANGUAGE_SERVER_TYPE" in \
    "typescript") \
        exec mcp-language-server \
            --mode="$MCP_MODE" \
            --port="$MCP_PORT" \
            --workspace="$WORKSPACE_DIR" \
            --lsp=typescript-language-server \
            -- --stdio ;; \
    "go") \
        exec mcp-language-server \
            --mode="$MCP_MODE" \
            --port="$MCP_PORT" \
            --workspace="$WORKSPACE_DIR" \
            --lsp=gopls ;; \
    "python") \
        exec mcp-language-server \
            --mode="$MCP_MODE" \
            --port="$MCP_PORT" \
            --workspace="$WORKSPACE_DIR" \
            --lsp=pyright-langserver \
            -- --stdio ;; \
    "rust") \
        exec mcp-language-server \
            --mode="$MCP_MODE" \
            --port="$MCP_PORT" \
            --workspace="$WORKSPACE_DIR" \
            --lsp=rust-analyzer ;; \
    "clangd") \
        exec mcp-language-server \
            --mode="$MCP_MODE" \
            --port="$MCP_PORT" \
            --workspace="$WORKSPACE_DIR" \
            --lsp=clangd ;; \
    *) \
        echo "Unsupported LANGUAGE_SERVER_TYPE: $LANGUAGE_SERVER_TYPE" >&2; \
        echo "Supported types: typescript, go, python, rust, clangd" >&2; \
        exit 1 ;; \
    esac

# Specialized single-language images for smaller size
FROM runtime-base AS typescript-only
COPY --from=node-base /usr/local/lib/node_modules /usr/local/lib/node_modules
COPY --from=node-base /usr/local/bin/typescript-language-server /usr/local/bin/
COPY --from=builder /app/mcp-language-server /usr/local/bin/
RUN mkdir -p /workspace
WORKDIR /workspace
EXPOSE 8080
CMD ["mcp-language-server", "--mode=http", "--port=8080", "--workspace=/workspace", "--lsp=typescript-language-server", "--", "--stdio"]

FROM runtime-base AS go-only
COPY --from=go-base /go/bin/gopls /usr/local/bin/
COPY --from=builder /app/mcp-language-server /usr/local/bin/
RUN mkdir -p /workspace
WORKDIR /workspace
EXPOSE 8080
CMD ["mcp-language-server", "--mode=http", "--port=8080", "--workspace=/workspace", "--lsp=gopls"]

FROM runtime-base AS python-only
COPY --from=node-base /usr/local/lib/node_modules/pyright /usr/local/lib/node_modules/pyright
COPY --from=node-base /usr/local/bin/pyright-langserver /usr/local/bin/
COPY --from=builder /app/mcp-language-server /usr/local/bin/
RUN mkdir -p /workspace
WORKDIR /workspace
EXPOSE 8080
CMD ["mcp-language-server", "--mode=http", "--port=8080", "--workspace=/workspace", "--lsp=pyright-langserver", "--", "--stdio"]

FROM runtime-base AS rust-only
COPY --from=rust-base /usr/local/rustup/toolchains/*/bin/rust-analyzer /usr/local/bin/
COPY --from=builder /app/mcp-language-server /usr/local/bin/
RUN mkdir -p /workspace
WORKDIR /workspace
EXPOSE 8080
CMD ["mcp-language-server", "--mode=http", "--port=8080", "--workspace=/workspace", "--lsp=rust-analyzer"]

FROM runtime-base AS clangd-only
COPY --from=clangd-base /usr/bin/clangd /usr/local/bin/
COPY --from=builder /app/mcp-language-server /usr/local/bin/
RUN mkdir -p /workspace
WORKDIR /workspace
EXPOSE 8080
CMD ["mcp-language-server", "--mode=http", "--port=8080", "--workspace=/workspace", "--lsp=clangd"]
