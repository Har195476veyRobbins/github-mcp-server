// Package main is the entry point for the GitHub MCP Server.
// It initializes and starts the Model Context Protocol server that exposes
// GitHub API functionality as tools for AI assistants.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time via ldflags.
	Version = "dev"
	// Commit is the git commit hash set at build time.
	Commit = "none"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	var (
		token   string
		transport string
		port    int
	)

	cmd := &cobra.Command{
		Use:     "github-mcp-server",
		Short:   "GitHub MCP Server - expose GitHub API as MCP tools",
		Version: fmt.Sprintf("%s (%s)", Version, Commit),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(cmd.Context(), token, transport, port)
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "GitHub personal access token (or set GITHUB_TOKEN env var)")
	cmd.Flags().StringVar(&transport, "transport", "stdio", "Transport type: stdio or sse")
	cmd.Flags().IntVar(&port, "port", 8080, "Port to listen on when using SSE transport")

	return cmd
}

func runServer(ctx context.Context, token, transport string, port int) error {
	// Resolve token from flag or environment variable.
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		return fmt.Errorf("GitHub token is required: set --token flag or GITHUB_TOKEN environment variable")
	}

	// Create the MCP server with GitHub tools.
	s, err := newGitHubMCPServer(token)
	if err != nil {
		return fmt.Errorf("failed to create MCP server: %w", err)
	}

	// Handle graceful shutdown.
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	switch transport {
	case "stdio":
		fmt.Fprintln(os.Stderr, "Starting GitHub MCP Server on stdio transport")
		return server.ServeStdio(s)
	case "sse":
		addr := fmt.Sprintf(":%d", port)
		fmt.Fprintf(os.Stderr, "Starting GitHub MCP Server on SSE transport at %s\n", addr)
		sseServer := server.NewSSEServer(s, server.WithBaseURL(fmt.Sprintf("http://localhost%s", addr)))
		return sseServer.Start(addr)
	default:
		return fmt.Errorf("unknown transport %q: must be 'stdio' or 'sse'", transport)
	}
}
