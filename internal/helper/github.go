package helper

import (
	"context"
	"os"

	"github.com/google/go-github/v49/github"
	"golang.org/x/oauth2"
)

func NewGitHubClent(ctx context.Context) *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
