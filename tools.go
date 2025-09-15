package main

import (
	"context"
	"fmt"

	"github.com/isaacphi/mcp-language-server/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool parameter types
type EditFileParams struct {
	FilePath string      `json:"filePath" jsonschema:"description=Path to the file to edit"`
	Edits    []TextEdit  `json:"edits" jsonschema:"description=List of edits to apply"`
}

type TextEdit struct {
	StartLine int    `json:"startLine" jsonschema:"description=Start line to replace inclusive one-indexed"`
	EndLine   int    `json:"endLine" jsonschema:"description=End line to replace inclusive one-indexed"`
	NewText   string `json:"newText" jsonschema:"description=Replacement text. Replace with the new text. Leave blank to remove lines."`
}

type DefinitionParams struct {
	SymbolName string `json:"symbolName" jsonschema:"description=The name of the symbol whose definition you want to find (e.g. 'mypackage.MyFunction', 'MyType.MyMethod')"`
}

type ReferencesParams struct {
	SymbolName string `json:"symbolName" jsonschema:"description=The name of the symbol to search for (e.g. 'mypackage.MyFunction', 'MyType')"`
}

type DiagnosticsParams struct {
	FilePath        string `json:"filePath" jsonschema:"description=The path to the file to get diagnostics for"`
	ContextLines    int    `json:"contextLines,omitempty" jsonschema:"description=Lines to include around each diagnostic,default=5"`
	ShowLineNumbers bool   `json:"showLineNumbers,omitempty" jsonschema:"description=If true adds line numbers to the output,default=true"`
}

type HoverParams struct {
	FilePath string `json:"filePath" jsonschema:"description=The path to the file to get hover information for"`
	Line     int    `json:"line" jsonschema:"description=The line number where the hover is requested (1-indexed)"`
	Column   int    `json:"column" jsonschema:"description=The column number where the hover is requested (1-indexed)"`
}

type RenameSymbolParams struct {
	FilePath string `json:"filePath" jsonschema:"description=The path to the file containing the symbol to rename"`
	Line     int    `json:"line" jsonschema:"description=The line number where the symbol is located (1-indexed)"`
	Column   int    `json:"column" jsonschema:"description=The column number where the symbol is located (1-indexed)"`
	NewName  string `json:"newName" jsonschema:"description=The new name for the symbol"`
}

func (s *mcpServer) registerTools() error {
	coreLogger.Debug("Registering MCP tools")

	// Edit file tool
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "edit_file",
		Description: "Apply multiple text edits to a file.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params EditFileParams) (*mcp.CallToolResult, any, error) {
		var edits []tools.TextEdit
		for _, edit := range params.Edits {
			edits = append(edits, tools.TextEdit{
				StartLine: edit.StartLine,
				EndLine:   edit.EndLine,
				NewText:   edit.NewText,
			})
		}

		coreLogger.Debug("Executing edit_file for file: %s", params.FilePath)
		response, err := tools.ApplyTextEdits(s.ctx, s.lspClient, params.FilePath, edits)
		if err != nil {
			coreLogger.Error("Failed to apply edits: %v", err)
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("failed to apply edits: %v", err)}},
				IsError: true,
			}, nil, nil
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: response}},
		}, nil, nil
	})

	// Definition tool
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "definition",
		Description: "Read the source code definition of a symbol (function, type, constant, etc.) from the codebase. Returns the complete implementation code where the symbol is defined.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params DefinitionParams) (*mcp.CallToolResult, any, error) {
		coreLogger.Debug("Executing definition for symbol: %s", params.SymbolName)
		text, err := tools.ReadDefinition(s.ctx, s.lspClient, params.SymbolName)
		if err != nil {
			coreLogger.Error("Failed to get definition: %v", err)
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("failed to get definition: %v", err)}},
				IsError: true,
			}, nil, nil
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: text}},
		}, nil, nil
	})

	// References tool
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "references",
		Description: "Find all usages and references of a symbol throughout the codebase. Returns a list of all files and locations where the symbol appears.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params ReferencesParams) (*mcp.CallToolResult, any, error) {
		coreLogger.Debug("Executing references for symbol: %s", params.SymbolName)
		text, err := tools.FindReferences(s.ctx, s.lspClient, params.SymbolName)
		if err != nil {
			coreLogger.Error("Failed to find references: %v", err)
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("failed to find references: %v", err)}},
				IsError: true,
			}, nil, nil
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: text}},
		}, nil, nil
	})

	// Diagnostics tool
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "diagnostics",
		Description: "Get diagnostic information for a specific file from the language server.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params DiagnosticsParams) (*mcp.CallToolResult, any, error) {
		contextLines := params.ContextLines
		if contextLines == 0 {
			contextLines = 5 // default value
		}
		showLineNumbers := params.ShowLineNumbers
		if params.ShowLineNumbers == false && contextLines != 0 {
			showLineNumbers = true // default value
		}

		coreLogger.Debug("Executing diagnostics for file: %s", params.FilePath)
		text, err := tools.GetDiagnosticsForFile(s.ctx, s.lspClient, params.FilePath, contextLines, showLineNumbers)
		if err != nil {
			coreLogger.Error("Failed to get diagnostics: %v", err)
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("failed to get diagnostics: %v", err)}},
				IsError: true,
			}, nil, nil
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: text}},
		}, nil, nil
	})

	// Hover tool
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "hover",
		Description: "Get hover information (type, documentation) for a symbol at the specified position.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params HoverParams) (*mcp.CallToolResult, any, error) {
		coreLogger.Debug("Executing hover for file: %s line: %d column: %d", params.FilePath, params.Line, params.Column)
		text, err := tools.GetHoverInfo(s.ctx, s.lspClient, params.FilePath, params.Line, params.Column)
		if err != nil {
			coreLogger.Error("Failed to get hover information: %v", err)
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("failed to get hover information: %v", err)}},
				IsError: true,
			}, nil, nil
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: text}},
		}, nil, nil
	})

	// Rename symbol tool
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "rename_symbol",
		Description: "Rename a symbol (variable, function, class, etc.) at the specified position and update all references throughout the codebase.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params RenameSymbolParams) (*mcp.CallToolResult, any, error) {
		coreLogger.Debug("Executing rename_symbol for file: %s line: %d column: %d newName: %s", params.FilePath, params.Line, params.Column, params.NewName)
		text, err := tools.RenameSymbol(s.ctx, s.lspClient, params.FilePath, params.Line, params.Column, params.NewName)
		if err != nil {
			coreLogger.Error("Failed to rename symbol: %v", err)
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("failed to rename symbol: %v", err)}},
				IsError: true,
			}, nil, nil
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: text}},
		}, nil, nil
	})

	coreLogger.Info("Successfully registered all MCP tools")
	return nil
}
