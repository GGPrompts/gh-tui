package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// helpers.go - Utility Functions
// Purpose: Common helper functions used across views

// truncateString truncates text to fit within maxWidth, adding "..." if needed
func truncateString(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if len(s) <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		return s[:maxWidth]
	}
	return s[:maxWidth-3] + "..."
}

// padRight pads a string with spaces to reach the target width
func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// formatTime formats timestamps in a user-friendly way
func formatTime(t time.Time) string {
	return t.Format("Jan 2, 15:04")
}

// formatTimeAgo shows relative time (e.g., "2h ago", "3d ago")
func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration.Seconds() < 60 {
		return "just now"
	}

	if duration.Minutes() < 60 {
		mins := int(duration.Minutes())
		return fmt.Sprintf("%dm ago", mins)
	}

	if duration.Hours() < 24 {
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh ago", hours)
	}

	days := int(duration.Hours() / 24)
	if days < 30 {
		return fmt.Sprintf("%dd ago", days)
	}

	months := days / 30
	if months < 12 {
		return fmt.Sprintf("%dmo ago", months)
	}

	years := months / 12
	return fmt.Sprintf("%dy ago", years)
}

// formatLabels formats a list of labels as a comma-separated string
func formatLabels(labels []Label) string {
	if len(labels) == 0 {
		return "none"
	}

	names := make([]string, len(labels))
	for i, label := range labels {
		names[i] = label.Name
	}
	return strings.Join(names, ", ")
}

// formatLanguage formats a language name, handling nil cases
func formatLanguage(lang *Language) string {
	if lang == nil {
		return "none"
	}
	return lang.Name
}

// formatVisibility formats visibility string with an icon
func formatVisibility(vis string) string {
	switch vis {
	case "PUBLIC":
		return "🌐 Public"
	case "PRIVATE":
		return "🔒 Private"
	default:
		return vis
	}
}

// formatStatus formats workflow run status with color/icon
func formatStatus(status, conclusion string) string {
	if status == "completed" {
		switch conclusion {
		case "success":
			return "✓ Success"
		case "failure":
			return "✗ Failure"
		case "cancelled":
			return "⊘ Cancelled"
		case "skipped":
			return "⊙ Skipped"
		default:
			return conclusion
		}
	}
	return status
}

// formatPRState formats PR state with an icon
func formatPRState(state string, isDraft bool) string {
	if isDraft {
		return "📝 Draft"
	}
	switch state {
	case "OPEN":
		return "🟢 Open"
	case "CLOSED":
		return "🔴 Closed"
	case "MERGED":
		return "🟣 Merged"
	default:
		return state
	}
}

// formatIssueState formats issue state with an icon
func formatIssueState(state string) string {
	switch state {
	case "OPEN":
		return "🟢 Open"
	case "CLOSED":
		return "🟣 Closed"
	default:
		return state
	}
}

// formatNumber formats numbers with k/M suffixes for large values
func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(n)/1000000)
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// makeHyperlink creates a clickable terminal hyperlink using OSC 8
// The displayed text can be truncated while the actual URL remains full
func makeHyperlink(url string, displayText string) string {
	// OSC 8 format: \x1b]8;;URL\x1b\\DISPLAY_TEXT\x1b]8;;\x1b\\
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, displayText)
}

// wrapText wraps text to fit within a given width
func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len() == 0 {
			currentLine.WriteString(word)
		} else if currentLine.Len()+1+len(word) <= width {
			currentLine.WriteString(" " + word)
		} else {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// openInBrowser opens a GitHub item in the default browser using gh CLI
// itemType: "pr", "issue", "repo", "run", "gist"
// identifier: PR number, issue number, repo name, run ID, or gist ID
// repo: repository path (e.g., "owner/repo") - optional for some types
func openInBrowser(itemType, identifier, repo string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd

		switch itemType {
		case "pr":
			// gh pr view <number> --web [--repo <repo>]
			if repo != "" {
				cmd = exec.Command("gh", "pr", "view", identifier, "--web", "--repo", repo)
			} else {
				cmd = exec.Command("gh", "pr", "view", identifier, "--web")
			}

		case "issue":
			// gh issue view <number> --web [--repo <repo>]
			if repo != "" {
				cmd = exec.Command("gh", "issue", "view", identifier, "--web", "--repo", repo)
			} else {
				cmd = exec.Command("gh", "issue", "view", identifier, "--web")
			}

		case "repo":
			// gh repo view <repo> --web
			cmd = exec.Command("gh", "repo", "view", identifier, "--web")

		case "run":
			// gh run view <id> --web [--repo <repo>]
			if repo != "" {
				cmd = exec.Command("gh", "run", "view", identifier, "--web", "--repo", repo)
			} else {
				cmd = exec.Command("gh", "run", "view", identifier, "--web")
			}

		case "gist":
			// gh gist view <id> --web
			cmd = exec.Command("gh", "gist", "view", identifier, "--web")

		default:
			return errMsg{err: fmt.Errorf("unknown item type: %s", itemType)}
		}

		if err := cmd.Run(); err != nil {
			return errMsg{err: fmt.Errorf("failed to open in browser: %w", err)}
		}

		return nil
	}
}

// createNewIssue opens the browser to create a new issue
func createNewIssue() tea.Cmd {
	return func() tea.Msg {
		// gh issue create --web opens browser with new issue form
		cmd := exec.Command("gh", "issue", "create", "--web")

		if err := cmd.Run(); err != nil {
			return errMsg{err: fmt.Errorf("failed to create issue: %w", err)}
		}

		return nil
	}
}

// toggleRepoStar stars or unstars a repository
func toggleRepoStar(repoNameWithOwner string) tea.Cmd {
	return func() tea.Msg {
		// First check if already starred
		checkCmd := exec.Command("gh", "repo", "view", repoNameWithOwner, "--json", "viewerHasStarred", "-q", ".viewerHasStarred")
		output, err := checkCmd.Output()
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to check star status: %w", err)}
		}

		isStarred := strings.TrimSpace(string(output)) == "true"

		var cmd *exec.Cmd
		if isStarred {
			// Unstar
			cmd = exec.Command("gh", "repo", "unstar", repoNameWithOwner)
		} else {
			// Star
			cmd = exec.Command("gh", "repo", "star", repoNameWithOwner)
		}

		if err := cmd.Run(); err != nil {
			return errMsg{err: fmt.Errorf("failed to toggle star: %w", err)}
		}

		// Return success message
		if isStarred {
			return statusMsg{message: "Repository unstarred"}
		}
		return statusMsg{message: "Repository starred ⭐"}
	}
}

// cloneRepository clones a repository to the current directory
func cloneRepository(repoNameWithOwner string) tea.Cmd {
	return func() tea.Msg {
		// gh repo clone opens in current directory
		cmd := exec.Command("gh", "repo", "clone", repoNameWithOwner)

		if err := cmd.Run(); err != nil {
			return errMsg{err: fmt.Errorf("failed to clone repository: %w", err)}
		}

		return statusMsg{message: fmt.Sprintf("Cloned %s successfully", repoNameWithOwner)}
	}
}

// viewPRDiff shows the diff for a pull request in the pager
func viewPRDiff(prNumber string) tea.Cmd {
	return func() tea.Msg {
		// gh pr diff shows the diff in the default pager (less, more, etc.)
		cmd := exec.Command("gh", "pr", "diff", prNumber)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return errMsg{err: fmt.Errorf("failed to view diff: %w", err)}
		}

		return nil
	}
}

// forkRepository forks a repository to the authenticated user's account
func forkRepository(repoNameWithOwner string) tea.Cmd {
	return func() tea.Msg {
		// gh repo fork creates a fork in the authenticated user's account
		cmd := exec.Command("gh", "repo", "fork", repoNameWithOwner, "--remote=false")

		if err := cmd.Run(); err != nil {
			return errMsg{err: fmt.Errorf("failed to fork repository: %w", err)}
		}

		return statusMsg{message: fmt.Sprintf("Forked %s successfully 🍴", repoNameWithOwner)}
	}
}

// closeIssue closes an issue
func closeIssue(issueNumber int) tea.Cmd {
	return func() tea.Msg {
		// gh issue close <number> closes the issue
		cmd := exec.Command("gh", "issue", "close", fmt.Sprintf("%d", issueNumber))

		if err := cmd.Run(); err != nil {
			return errMsg{err: fmt.Errorf("failed to close issue: %w", err)}
		}

		return statusMsg{message: fmt.Sprintf("Issue #%d closed", issueNumber)}
	}
}

// reopenIssue reopens a closed issue
func reopenIssue(issueNumber int) tea.Cmd {
	return func() tea.Msg {
		// gh issue reopen <number> reopens the issue
		cmd := exec.Command("gh", "issue", "reopen", fmt.Sprintf("%d", issueNumber))

		if err := cmd.Run(); err != nil {
			return errMsg{err: fmt.Errorf("failed to reopen issue: %w", err)}
		}

		return statusMsg{message: fmt.Sprintf("Issue #%d reopened", issueNumber)}
	}
}

// viewWorkflowLogs shows the logs for a workflow run in the pager
func viewWorkflowLogs(runId string) tea.Cmd {
	return func() tea.Msg {
		// gh run view shows the logs in the default pager (less, more, etc.)
		cmd := exec.Command("gh", "run", "view", runId, "--log")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return errMsg{err: fmt.Errorf("failed to view logs: %w", err)}
		}

		return nil
	}
}

