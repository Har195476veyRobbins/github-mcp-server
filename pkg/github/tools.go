package github

import (
	"context"
	"fmt"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/google/go-github/v67/github"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registers all GitHub MCP tools with the given server.
// Each tool corresponds to a GitHub API operation exposed via the MCP protocol.
func RegisterTools(s *server.MCPServer, client *github.Client, t translations.TranslationHelperFunc) {
	registerRepositoryTools(s, client, t)
	registerIssueTools(s, client, t)
	registerPullRequestTools(s, client, t)
}

// registerRepositoryTools registers tools related to GitHub repositories.
func registerRepositoryTools(s *server.MCPServer, client *github.Client, t translations.TranslationHelperFunc) {
	s.AddTool(
		mcp.NewTool(
			"get_repository",
			mcp.WithDescription(t("TOOL_GET_REPOSITORY_DESCRIPTION", "Get details about a GitHub repository")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("The account owner of the repository"),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("The name of the repository"),
			),
		),
		getRepositoryHandler(client),
	)
}

// registerIssueTools registers tools related to GitHub issues.
func registerIssueTools(s *server.MCPServer, client *github.Client, t translations.TranslationHelperFunc) {
	s.AddTool(
		mcp.NewTool(
			"list_issues",
			mcp.WithDescription(t("TOOL_LIST_ISSUES_DESCRIPTION", "List issues for a GitHub repository")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("The account owner of the repository"),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("The name of the repository"),
			),
			mcp.WithString("state",
				// Personal preference: default to "all" so we see both open and closed issues
				mcp.Description("Filter issues by state: open, closed, or all (default: all)"),
			),
		),
		listIssuesHandler(client),
	)
}

// registerPullRequestTools registers tools related to GitHub pull requests.
func registerPullRequestTools(s *server.MCPServer, client *github.Client, t translations.TranslationHelperFunc) {
	s.AddTool(
		mcp.NewTool(
			"list_pull_requests",
			mcp.WithDescription(t("TOOL_LIST_PULL_REQUESTS_DESCRIPTION", "List pull requests for a GitHub repository")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("The account owner of the repository"),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("The name of the repository"),
			),
			mcp.WithString("state",
				// Personal preference: default to "all" so we see both open and closed PRs
				mcp.Description("Filter pull requests by state: open, closed, or all (default: all)"),
			),
		),
		listPullRequestsHandler(client),
	)
}

// getRepositoryHandler returns an MCP tool handler that fetches a single repository.
func getRepositoryHandler(client *github.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		owner, err := req.RequireString("owner")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		repo, err := r
