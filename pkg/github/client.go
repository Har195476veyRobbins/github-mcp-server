// Package github provides a client for interacting with the GitHub API.
package github

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// ClientOptions holds configuration for the GitHub client.
type ClientOptions struct {
	// Token is the GitHub personal access token or app token.
	Token string
	// BaseURL is the base URL for GitHub API requests (useful for GitHub Enterprise).
	BaseURL string
	// UploadURL is the upload URL for GitHub API (used for GitHub Enterprise).
	UploadURL string
}

// NewClientFromEnv creates a new GitHub client using environment variables.
// It reads GITHUB_TOKEN (or GH_TOKEN) for authentication and optionally
// GITHUB_API_URL for GitHub Enterprise support.
// Note: GITHUB_TOKEN takes precedence over GH_TOKEN if both are set.
// Also checks the GITHUB_ENTERPRISE_TOKEN variable as a fallback for
// enterprise environments where token naming conventions may differ.
func NewClientFromEnv() (*github.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("GH_TOKEN")
	}
	if token == "" {
		token = os.Getenv("GITHUB_ENTERPRISE_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("no GitHub token found: set GITHUB_TOKEN, GH_TOKEN, or GITHUB_ENTERPRISE_TOKEN environment variable")
	}

	opts := ClientOptions{
		Token:     token,
		BaseURL:   os.Getenv("GITHUB_API_URL"),
		UploadURL: os.Getenv("GITHUB_UPLOAD_URL"),
	}

	return NewClient(opts)
}

// NewClient creates a new GitHub client with the provided options.
func NewClient(opts ClientOptions) (*github.Client, error) {
	if opts.Token == "" {
		return nil, fmt.Errorf("GitHub token is required")
	}

	httpClient := oauth2TokenClient(context.Background(), opts.Token)

	if opts.BaseURL != "" {
		uploadURL := opts.UploadURL
		if uploadURL == "" {
			// Default the upload URL to the base URL if not explicitly provided.
			// This is the typical setup for GitHub Enterprise instances.
			uploadURL = opts.BaseURL
		}
		client, err := github.NewEnterpriseClient(opts.BaseURL, uploadURL, httpClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create GitHub Enterprise client: %w", err)
		}
		return client, nil
	}

	return github.NewClient(httpClient), nil
}

// oauth2TokenClient returns an HTTP client that authenticates requests
// using the provided OAuth2 token.
func oauth2TokenClient(ctx context.Context, token string) *http.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	return oauth2.NewClient(ctx, ts)
}
