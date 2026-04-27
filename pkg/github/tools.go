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
				mcp.Description("Filter issues by state: open, closed, or all (default: open)"),
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
				mcp.Description("Filter pull requests by state: open, closed, or all (default: open)"),
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
		repo, err := req.RequireString("repo")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		repository, _, err := client.Repositories.Get(ctx, owner, repo)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get repository: %v", err)), nil
		}

		return mcp.NewToolResultText(repository.GetFullName()), nil
	}
}

// listIssuesHandler returns an MCP tool handler that lists issues for a repository.
func listIssuesHandler(client *github.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		owner, err := req.RequireString("owner")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		repo, err := req.RequireString("repo")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		state := req.GetString("state", "open")

		issues, _, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
			State: state,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to list issues: %v", err)), nil
		}

		var result string
		for _, issue := range issues {
			result += fmt.Sprintf("#%d: %s\n", issue.GetNumber(), issue.GetTitle())
		}
		return mcp.NewToolResultText(result), nil
	}
}

// listPullRequestsHandler returns an MCP tool handler that lists pull requests for a repository.
func listPullRequestsHandler(client *github.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		owner, err := req.RequireString("owner")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		repo, err := req.RequireString("repo")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		state := req.GetString("state", "open")

		prs, _, err := client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
			State: state,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to list pull requests: %v", err)), nil
		}

		var result string
		for _, pr := range prs {
			result += fmt.Sprintf("#%d: %s\n", pr.GetNumber(), pr.GetTitle())
		}
		return mcp.NewToolResultText(result), nil
	}
}
