package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// github.go - GitHub CLI Integration
// Purpose: All GitHub CLI commands and data fetching
// Uses `gh` CLI as backend (must be authenticated)

// checkGitHubAuth verifies gh CLI is authenticated
func checkGitHubAuth() error {
	cmd := exec.Command("gh", "auth", "status")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gh not authenticated. Run: gh auth login")
	}
	return nil
}

// fetchPullRequests retrieves PRs using gh CLI
func fetchPullRequests(repo string) tea.Cmd {
	return func() tea.Msg {
		if repo == "" {
			// Try to get current repo
			cmd := exec.Command("gh", "repo", "view", "--json", "nameWithOwner")
			output, err := cmd.Output()
			if err != nil {
				return prLoadedMsg{err: fmt.Errorf("no repo specified and not in a git repo")}
			}
			var repoData struct {
				NameWithOwner string `json:"nameWithOwner"`
			}
			if err := json.Unmarshal(output, &repoData); err != nil {
				return prLoadedMsg{err: err}
			}
			repo = repoData.NameWithOwner
		}

		cmd := exec.Command("gh", "pr", "list",
			"--repo", repo,
			"--json", "number,title,state,author,createdAt,updatedAt,headRefName,baseRefName,isDraft,reviewDecision,mergeable,url",
			"--limit", "100")

		output, err := cmd.Output()
		if err != nil {
			return prLoadedMsg{err: fmt.Errorf("gh pr list failed: %w", err)}
		}

		var prs []PullRequest
		if err := json.Unmarshal(output, &prs); err != nil {
			return prLoadedMsg{err: fmt.Errorf("parse error: %w", err)}
		}

		return prLoadedMsg{prs: prs}
	}
}

// fetchIssues retrieves issues using gh CLI
func fetchIssues(repo string) tea.Cmd {
	return func() tea.Msg {
		if repo == "" {
			cmd := exec.Command("gh", "repo", "view", "--json", "nameWithOwner")
			output, err := cmd.Output()
			if err != nil {
				return issuesLoadedMsg{err: fmt.Errorf("no repo specified and not in a git repo")}
			}
			var repoData struct {
				NameWithOwner string `json:"nameWithOwner"`
			}
			if err := json.Unmarshal(output, &repoData); err != nil {
				return issuesLoadedMsg{err: err}
			}
			repo = repoData.NameWithOwner
		}

		cmd := exec.Command("gh", "issue", "list",
			"--repo", repo,
			"--json", "number,title,state,author,createdAt,updatedAt,labels,assignees,milestone,url",
			"--limit", "100")

		output, err := cmd.Output()
		if err != nil {
			return issuesLoadedMsg{err: fmt.Errorf("gh issue list failed: %w", err)}
		}

		var issues []Issue
		if err := json.Unmarshal(output, &issues); err != nil {
			return issuesLoadedMsg{err: fmt.Errorf("parse error: %w", err)}
		}

		return issuesLoadedMsg{issues: issues}
	}
}

// fetchRepositories retrieves repos using gh CLI
func fetchRepositories(owner string) tea.Cmd {
	return func() tea.Msg {
		if owner == "" {
			// Get current user
			cmd := exec.Command("gh", "api", "user", "--jq", ".login")
			output, err := cmd.Output()
			if err != nil {
				return reposLoadedMsg{err: fmt.Errorf("failed to get current user: %w", err)}
			}
			owner = string(output[:len(output)-1]) // Remove newline
		}

		cmd := exec.Command("gh", "repo", "list", owner,
			"--json", "name,nameWithOwner,description,stargazerCount,forkCount,primaryLanguage,visibility,url",
			"--limit", "100")

		output, err := cmd.Output()
		if err != nil {
			return reposLoadedMsg{err: fmt.Errorf("gh repo list failed: %w", err)}
		}

		var repos []Repository
		if err := json.Unmarshal(output, &repos); err != nil {
			return reposLoadedMsg{err: fmt.Errorf("parse error: %w", err)}
		}

		return reposLoadedMsg{repos: repos}
	}
}

// fetchWorkflowRuns retrieves workflow runs using gh CLI
func fetchWorkflowRuns(repo string) tea.Cmd {
	return func() tea.Msg {
		if repo == "" {
			cmd := exec.Command("gh", "repo", "view", "--json", "nameWithOwner")
			output, err := cmd.Output()
			if err != nil {
				return workflowsLoadedMsg{err: fmt.Errorf("no repo specified and not in a git repo")}
			}
			var repoData struct {
				NameWithOwner string `json:"nameWithOwner"`
			}
			if err := json.Unmarshal(output, &repoData); err != nil {
				return workflowsLoadedMsg{err: err}
			}
			repo = repoData.NameWithOwner
		}

		cmd := exec.Command("gh", "run", "list",
			"--repo", repo,
			"--json", "databaseId,name,status,conclusion,headBranch,headSha,number,createdAt,url",
			"--limit", "50")

		output, err := cmd.Output()
		if err != nil {
			return workflowsLoadedMsg{err: fmt.Errorf("gh run list failed: %w", err)}
		}

		var runs []WorkflowRun
		if err := json.Unmarshal(output, &runs); err != nil {
			return workflowsLoadedMsg{err: fmt.Errorf("parse error: %w", err)}
		}

		return workflowsLoadedMsg{runs: runs}
	}
}

// fetchGists retrieves gists using gh CLI API
func fetchGists() tea.Cmd {
	return func() tea.Msg {
		// gh gist list doesn't support --json, so we use the API directly
		cmd := exec.Command("gh", "api", "/gists", "--paginate", "-q", ".[0:50]")

		output, err := cmd.Output()
		if err != nil {
			return gistsLoadedMsg{err: fmt.Errorf("gh api /gists failed: %w", err)}
		}

		// Parse the API response which has different field names
		var apiGists []struct {
			ID          string                       `json:"id"`
			Description string                       `json:"description"`
			Public      bool                         `json:"public"`
			Files       map[string]map[string]string `json:"files"`
			CreatedAt   string                       `json:"created_at"`
			UpdatedAt   string                       `json:"updated_at"`
			HTMLURL     string                       `json:"html_url"`
		}

		if err := json.Unmarshal(output, &apiGists); err != nil {
			return gistsLoadedMsg{err: fmt.Errorf("parse error: %w", err)}
		}

		// Transform to our Gist structure
		gists := make([]Gist, len(apiGists))
		for i, apiGist := range apiGists {
			// Convert files map to array
			files := make([]GistFile, 0, len(apiGist.Files))
			for filename := range apiGist.Files {
				files = append(files, GistFile{Filename: filename})
			}

			// Parse timestamps
			createdAt, _ := time.Parse(time.RFC3339, apiGist.CreatedAt)
			updatedAt, _ := time.Parse(time.RFC3339, apiGist.UpdatedAt)

			gists[i] = Gist{
				ID:          apiGist.ID,
				Description: apiGist.Description,
				Public:      apiGist.Public,
				Files:       files,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				URL:         apiGist.HTMLURL,
			}
		}

		return gistsLoadedMsg{gists: gists}
	}
}
