package main

import (
	"encoding/json"
	"fmt"
	"os/exec"

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
			"--json", "databaseId,name,status,conclusion,headBranch,headSha,runNumber,createdAt,url",
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

// fetchGists retrieves gists using gh CLI
func fetchGists() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("gh", "gist", "list",
			"--json", "id,description,public,files,createdAt,updatedAt,url",
			"--limit", "50")

		output, err := cmd.Output()
		if err != nil {
			return gistsLoadedMsg{err: fmt.Errorf("gh gist list failed: %w", err)}
		}

		var gists []Gist
		if err := json.Unmarshal(output, &gists); err != nil {
			return gistsLoadedMsg{err: fmt.Errorf("parse error: %w", err)}
		}

		return gistsLoadedMsg{gists: gists}
	}
}
