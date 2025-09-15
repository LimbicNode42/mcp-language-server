package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/isaacphi/mcp-language-server/internal/logging"
	"github.com/isaacphi/mcp-language-server/internal/lsp"
	"github.com/isaacphi/mcp-language-server/internal/watcher"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Create a logger for the core component
var coreLogger = logging.NewLogger(logging.Core)

type config struct {
	workspaceDir string
	lspCommand   string
	lspArgs      []string
	mode         string // "stdio" or "http"
	port         int    // for HTTP mode
}

type mcpServer struct {
	config           config
	lspClient        *lsp.Client
	mcpServer        *mcp.Server
	ctx              context.Context
	cancelFunc       context.CancelFunc
	workspaceWatcher *watcher.WorkspaceWatcher
}

func parseConfig() (*config, error) {
	cfg := &config{}
	flag.StringVar(&cfg.workspaceDir, "workspace", "", "Path to workspace directory")
	flag.StringVar(&cfg.lspCommand, "lsp", "", "LSP command to run (args should be passed after --)")
	flag.StringVar(&cfg.mode, "mode", "stdio", "Transport mode: 'stdio' or 'http'")
	flag.IntVar(&cfg.port, "port", 8080, "Port for HTTP mode (ignored for stdio mode)")
	flag.Parse()

	// Get remaining args after -- as LSP arguments
	cfg.lspArgs = flag.Args()

	// Validate transport mode
	if cfg.mode != "stdio" && cfg.mode != "http" {
		return nil, fmt.Errorf("mode must be 'stdio' or 'http'")
	}

	// Validate workspace directory
	if cfg.workspaceDir == "" {
		return nil, fmt.Errorf("workspace directory is required")
	}

	workspaceDir, err := filepath.Abs(cfg.workspaceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for workspace: %v", err)
	}
	cfg.workspaceDir = workspaceDir

	if _, err := os.Stat(cfg.workspaceDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("workspace directory does not exist: %s", cfg.workspaceDir)
	}

	// Validate LSP command
	if cfg.lspCommand == "" {
		return nil, fmt.Errorf("LSP command is required")
	}

	if _, err := exec.LookPath(cfg.lspCommand); err != nil {
		return nil, fmt.Errorf("LSP command not found: %s", cfg.lspCommand)
	}

	return cfg, nil
}

func newServer(config *config) (*mcpServer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &mcpServer{
		config:     *config,
		ctx:        ctx,
		cancelFunc: cancel,
	}, nil
}

func (s *mcpServer) initializeLSP() error {
	coreLogger.Info("Changing to workspace directory: %s", s.config.workspaceDir)
	if err := os.Chdir(s.config.workspaceDir); err != nil {
		return fmt.Errorf("failed to change to workspace directory: %v", err)
	}

	coreLogger.Info("Creating LSP client with command: %s, args: %v", s.config.lspCommand, s.config.lspArgs)
	client, err := lsp.NewClient(s.config.lspCommand, s.config.lspArgs...)
	if err != nil {
		return fmt.Errorf("failed to create LSP client: %v", err)
	}
	s.lspClient = client
	s.workspaceWatcher = watcher.NewWorkspaceWatcher(client)

	coreLogger.Info("Initializing LSP client...")
	initResult, err := client.InitializeLSPClient(s.ctx, s.config.workspaceDir)
	if err != nil {
		return fmt.Errorf("initialize failed: %v", err)
	}

	coreLogger.Info("LSP server capabilities received")
	coreLogger.Debug("Server capabilities: %+v", initResult.Capabilities)

	coreLogger.Info("Starting workspace watcher...")
	go s.workspaceWatcher.WatchWorkspace(s.ctx, s.config.workspaceDir)
	
	coreLogger.Info("Waiting for LSP server to be ready...")
	err = client.WaitForServerReady(s.ctx)
	if err != nil {
		return fmt.Errorf("LSP server ready wait failed: %v", err)
	}
	coreLogger.Info("LSP server is ready")
	return nil
}

func (s *mcpServer) start() error {
	coreLogger.Info("Initializing LSP client with command: %s, args: %v", s.config.lspCommand, s.config.lspArgs)
	coreLogger.Info("Workspace directory: %s", s.config.workspaceDir)
	
	if err := s.initializeLSP(); err != nil {
		coreLogger.Error("LSP initialization failed: %v", err)
		return err
	}
	coreLogger.Info("LSP client initialized successfully")

	s.mcpServer = mcp.NewServer(&mcp.Implementation{
		Name:    "MCP Language Server",
		Version: "v0.0.2",
	}, nil)

	coreLogger.Info("Registering MCP tools...")
	err := s.registerTools()
	if err != nil {
		coreLogger.Error("Tool registration failed: %v", err)
		return fmt.Errorf("tool registration failed: %v", err)
	}
	coreLogger.Info("MCP tools registered successfully")

	switch s.config.mode {
	case "stdio":
		coreLogger.Info("Starting MCP server in stdio mode")
		return s.mcpServer.Run(s.ctx, &mcp.StdioTransport{})
	case "http":
		addr := fmt.Sprintf(":%d", s.config.port)
		coreLogger.Info("Starting MCP server in HTTP mode")
		coreLogger.Info("Server will bind to address: %s", addr)
		coreLogger.Info("Full server URL will be: http://0.0.0.0%s", addr)
		
		handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
			coreLogger.Info("HTTP request received: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			return s.mcpServer
		}, nil)
		
		// Add logging middleware
		loggedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			coreLogger.Info("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			handler.ServeHTTP(w, r)
			coreLogger.Info("Response completed for %s %s in %v", r.Method, r.URL.Path, time.Since(start))
		})
		
		httpServer := &http.Server{
			Addr:    addr,
			Handler: loggedHandler,
		}
		
		// Start server in a goroutine so we can handle shutdown
		go func() {
			<-s.ctx.Done()
			coreLogger.Info("Shutting down HTTP server")
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			httpServer.Shutdown(shutdownCtx)
		}()
		
		coreLogger.Info("About to call ListenAndServe() on %s", addr)
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			coreLogger.Error("HTTP server failed: %v", err)
			return fmt.Errorf("HTTP server failed: %v", err)
		}
		coreLogger.Info("HTTP server stopped")
		return nil
	default:
		return fmt.Errorf("unsupported mode: %s", s.config.mode)
	}
}

func main() {
	coreLogger.Info("MCP Language Server starting")
	coreLogger.Info("Process ID: %d, Parent ID: %d", os.Getpid(), os.Getppid())

	done := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	config, err := parseConfig()
	if err != nil {
		coreLogger.Fatal("Failed to parse config: %v", err)
	}

	coreLogger.Info("Configuration loaded:")
	coreLogger.Info("  Mode: %s", config.mode)
	coreLogger.Info("  Port: %d", config.port)
	coreLogger.Info("  Workspace: %s", config.workspaceDir)
	coreLogger.Info("  LSP Command: %s", config.lspCommand)
	coreLogger.Info("  LSP Args: %v", config.lspArgs)

	server, err := newServer(config)
	if err != nil {
		coreLogger.Fatal("Failed to create server: %v", err)
	}

	// Parent process monitoring channel
	parentDeath := make(chan struct{})

	// Monitor parent process termination
	// Claude desktop does not properly kill child processes for MCP servers
	go func() {
		ppid := os.Getppid()
		coreLogger.Debug("Monitoring parent process: %d", ppid)

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				currentPpid := os.Getppid()
				if currentPpid != ppid && (currentPpid == 1 || ppid == 1) {
					coreLogger.Info("Parent process %d terminated (current ppid: %d), initiating shutdown", ppid, currentPpid)
					close(parentDeath)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Handle shutdown triggers
	go func() {
		select {
		case sig := <-sigChan:
			coreLogger.Info("Received signal %v in PID: %d", sig, os.Getpid())
			cleanup(server, done)
		case <-parentDeath:
			coreLogger.Info("Parent death detected, initiating shutdown")
			cleanup(server, done)
		}
	}()

	if err := server.start(); err != nil {
		coreLogger.Error("Server error: %v", err)
		cleanup(server, done)
		os.Exit(1)
	}

	<-done
	coreLogger.Info("Server shutdown complete for PID: %d", os.Getpid())
	os.Exit(0)
}

func cleanup(s *mcpServer, done chan struct{}) {
	coreLogger.Info("Cleanup initiated for PID: %d", os.Getpid())

	// Create a context with timeout for shutdown operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.lspClient != nil {
		coreLogger.Info("Closing open files")
		s.lspClient.CloseAllFiles(ctx)

		// Create a shorter timeout context for the shutdown request
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 500*time.Millisecond)
		defer shutdownCancel()

		// Run shutdown in a goroutine with timeout to avoid blocking if LSP doesn't respond
		shutdownDone := make(chan struct{})
		go func() {
			coreLogger.Info("Sending shutdown request")
			if err := s.lspClient.Shutdown(shutdownCtx); err != nil {
				coreLogger.Error("Shutdown request failed: %v", err)
			}
			close(shutdownDone)
		}()

		// Wait for shutdown with timeout
		select {
		case <-shutdownDone:
			coreLogger.Info("Shutdown request completed")
		case <-time.After(1 * time.Second):
			coreLogger.Warn("Shutdown request timed out, proceeding with exit")
		}

		coreLogger.Info("Sending exit notification")
		if err := s.lspClient.Exit(ctx); err != nil {
			coreLogger.Error("Exit notification failed: %v", err)
		}

		coreLogger.Info("Closing LSP client")
		if err := s.lspClient.Close(); err != nil {
			coreLogger.Error("Failed to close LSP client: %v", err)
		}
	}

	// Send signal to the done channel
	select {
	case <-done: // Channel already closed
	default:
		close(done)
	}

	coreLogger.Info("Cleanup completed for PID: %d", os.Getpid())
}
