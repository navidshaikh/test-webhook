//package main
package github

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubCreds struct {
	Token string
}

type Github struct {
	Org    string
	Repo   string
	Creds  GithubCreds
	Client *github.Client
}

type GithubComment struct {
	Body   string
	Author string
	Id     string
}

type TriggerComment struct {
	PR        int
	CommentID string
	Sha       string
}

func DefaultGithub(ctx context.Context, org, repo string) (*Github, error) {
	var token string
	if v := os.Getenv("GH_ACCESS_TOKEN"); v != "" {
		token = v
	}
	if token == "" {
		return nil, fmt.Errorf("GH_ACCESS_TOKEN is not in the environment")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &Github{
		Org:    org,
		Repo:   repo,
		Creds:  GithubCreds{Token: token},
		Client: client,
	}, nil
}

/*
func (g *Github) GetOpenPRs(ctx context.Context, opt *github.PullRequestListOptions) ([]*github.PullRequest, error) {
	openPRs, _, err := g.Client.PullRequests.List(ctx, g.Org, g.Repo, opt)
	if err != nil {
		return nil, err
	}
	return openPRs, nil
}

func (g *Github) GetPRComments(ctx context.Context, pr int, opt *github.IssueListCommentsOptions) ([]*github.IssueComment, error) {
	comments, _, err := g.Client.Issues.ListComments(ctx, g.Org, g.Repo, pr, opt)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (g *Github) ParsePRComments(ctx context.Context, pattern string, prno int, opt *github.IssueListCommentsOptions) (*github.IssueComment, error) {
	comments, err := g.GetPRComments(ctx, opt, prno)
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		match, err := regexp.MatchString(pattern, comment.GetBody())
		if err != nil {
			return nil, err
		}
		if match {
			return comment, nil
		}
	}
	return nil, nil
}

// TODO
func (g *Github) IsTrustedReviewer(ctx context.Context, username string) (bool, error) {
	opt := &github.RepositoryContentGetOptions{}
	file, _, _, err := g.Client.Repositories.GetContents(ctx, "CODEOWNERS", opt)
	if err != nil {
		return false, err
	}
	content, _ := file.GetContent()
	teams_regex, err := regexp.Compile("@[[:alnum:]-]+/[[:alnum:]-]+")
	if err != nil {
		return err
	}
	users_regex, err := regexp.Compile("@[[:alnum:]-]+[^/]")
	if err != nil {
		return err
	}
	teams := teams_regex.FindAllStringSubmatch(content, -1)
	for _, team := range teams {
		team = strings.Split((strings.Split(team[0], "@"))[1], "/")
		team_org := team[0]
		team_name := team[1]
		fmt.Println(team_org)
		fmt.Println(team_name)

	}
	users := users_regex.FindAllStringSubmatch(content, -1)
	fmt.Printf("USERS")
	for _, user := range users {
		user = strings.Split(user[0], "@")
		fmt.Println(user)
	}
	// two types
	// teams -> @[azAZ-]/[azAZ-]
	// team_org -> teams[0]
	// team -> teams[1]
	// codeowners -> @[azAZ09]

	// if username in codeowners || team
	// 	return true
	// else
	// 	return false

	// client.Teams.ListTeams()
	return false, nil
}

func (g *Github) PostGithubComment(ctx context.Context, comment string, prno int) error {
	return g.Client.Issues.CreateComment(ctx, g.Org, g.Repo, prno, &github.IssueComment{Body: &comment})
}
*/

func (g *Github) GetLatestPRCommit(ctx context.Context, prno int) (string, error) {
	commits, _, err := g.Client.PullRequests.ListCommits(ctx, g.Org, g.Repo, prno, &github.ListOptions{})
	if err != nil {
		return "", nil
	}
	if commits != nil {
		return *commits[len(commits)-1].SHA, nil
	}
	return "", fmt.Errorf("failed to find the latest commit of PR %d", prno)
}

/*
func main() {
	ctx := context.Background()
	g, err := DefaultGithub(ctx, "navidshaikh", "test-webhook")
	if err != nil {
		fmt.Println(err)
	}
	commit, err := g.GetLatestPRCommit(ctx, 6)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(commit)
}
*/
