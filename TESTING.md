# Testing MCP Language Server with Example Project

## Quick Test

1. **Mount the example project as workspace**:
   ```bash
   # Update your .env file:
   echo "MCP_TYPESCRIPT_PORT=3006" > .env
   echo "WORKSPACE_PATH=./example-projects/typescript-sample" >> .env
   ```

2. **Start the TypeScript MCP service**:
   ```bash
   docker-compose up mcp-typescript
   ```

3. **Test the endpoints**:
   ```bash
   # List available tools
   curl -X POST http://localhost:3006/mcp \
     -H "Content-Type: application/json" \
     -d '{"method": "list_tools", "params": {}}'

   # Get file diagnostics
   curl -X POST http://localhost:3006/mcp \
     -H "Content-Type: application/json" \
     -d '{"method": "call_tool", "params": {"name": "get_diagnostics", "arguments": {"uri": "file:///workspace/src/index.ts"}}}'

   # Get hover information
   curl -X POST http://localhost:3006/mcp \
     -H "Content-Type: application/json" \
     -d '{"method": "call_tool", "params": {"name": "get_hover_info", "arguments": {"uri": "file:///workspace/src/index.ts", "line": 5, "character": 10}}}'
   ```

## Expected Behavior

With a real TypeScript project mounted:
- Container should start successfully (no hanging)
- TypeScript language server will analyze the project files
- Diagnostics should show actual TypeScript errors
- Hover info should work on symbols
- Definition jumping should work between files

## Current Limitations

- **Single Project Only**: Each container instance serves exactly one project
- **Pre-mounted Workspace**: Project must exist before container starts
- **Language-Specific**: Need separate containers for different languages
- **No Dynamic Projects**: Cannot add/switch projects without restart

## Future Multi-Project Vision

```bash
# Future API (not implemented yet)
curl -X POST http://localhost:3006/projects \
  -d '{"name": "my-app", "path": "/path/to/project", "language": "typescript"}'

curl -X GET http://localhost:3006/projects
# -> [{"id": "proj-1", "name": "my-app", "language": "typescript", "active": true}]

curl -X POST http://localhost:3006/mcp \
  -d '{"method": "call_tool", "params": {"name": "get_diagnostics", "arguments": {"project_id": "proj-1", "uri": "src/index.ts"}}}'
```
