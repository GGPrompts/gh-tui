package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// checkMicroAvailable checks if micro editor is installed
func checkMicroAvailable() bool {
	_, err := exec.LookPath("micro")
	return err == nil
}

// downloadGistToTemp downloads a gist file to a temporary file
// For multi-file gists, downloads the first file
func downloadGistToTemp(gistID string, filename string) (string, error) {
	// Create temp file with gist ID in the name for easy identification
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("gh-tui-gist-%s-%s", gistID, filename))

	// Use gh CLI to download gist content
	cmd := exec.Command("gh", "gist", "view", gistID, "--filename", filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to download gist %s (file: %s): %s - %w", gistID, filename, string(output), err)
	}

	// Write to temp file
	if err := os.WriteFile(tempFile, output, 0644); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}

	return tempFile, nil
}

// calculateFileHash calculates SHA256 hash of a file
func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// openGistInMicro opens a gist file in micro editor
// readonly: if true, opens in read-only mode
func openGistInMicro(gistID string, filename string, readonly bool) tea.Cmd {
	// Download gist to temp file first
	tempFile, err := downloadGistToTemp(gistID, filename)
	if err != nil {
		return func() tea.Msg {
			return gistEditorFinishedMsg{
				gistID: gistID,
				err:    err,
			}
		}
	}

	// Calculate hash before editing (to detect modifications)
	hashBefore, err := calculateFileHash(tempFile)
	if err != nil {
		os.Remove(tempFile)
		return func() tea.Msg {
			return gistEditorFinishedMsg{
				gistID:       gistID,
				tempFilePath: tempFile,
				err:          fmt.Errorf("failed to calculate file hash: %w", err),
			}
		}
	}

	// Prepare micro command
	var cmd *exec.Cmd
	if readonly {
		cmd = exec.Command("micro", "-readonly", "true", tempFile)
	} else {
		cmd = exec.Command("micro", tempFile)
	}

	// Return a sequence that clears screen then opens editor
	return tea.Sequence(
		tea.ClearScreen,
		tea.ExecProcess(cmd, func(err error) tea.Msg {
			// After editor exits, check if file was modified
			wasModified := false
			if !readonly && err == nil {
				hashAfter, hashErr := calculateFileHash(tempFile)
				if hashErr == nil {
					wasModified = (hashBefore != hashAfter)
				}
			}

			return gistEditorFinishedMsg{
				gistID:       gistID,
				tempFilePath: tempFile,
				wasModified:  wasModified,
				isNewGist:    false,
				err:          err,
			}
		}),
	)
}

// createNewGistInMicro opens micro with an empty file for creating a new gist
func createNewGistInMicro() tea.Cmd {
	// Create temp file with placeholder name
	tempFile := filepath.Join(os.TempDir(), "gh-tui-new-gist.txt")

	// Create empty file
	if err := os.WriteFile(tempFile, []byte(""), 0644); err != nil {
		return func() tea.Msg {
			return gistEditorFinishedMsg{
				isNewGist: true,
				err:       fmt.Errorf("failed to create temp file: %w", err),
			}
		}
	}

	// Open in micro
	cmd := exec.Command("micro", tempFile)

	// Return a sequence that clears screen then opens editor
	return tea.Sequence(
		tea.ClearScreen,
		tea.ExecProcess(cmd, func(err error) tea.Msg {
			// Check if user wrote anything
			content, readErr := os.ReadFile(tempFile)
			wasModified := readErr == nil && len(strings.TrimSpace(string(content))) > 0

			return gistEditorFinishedMsg{
				tempFilePath: tempFile,
				wasModified:  wasModified,
				isNewGist:    true,
				err:          err,
			}
		}),
	)
}

// uploadGistChanges uploads modified gist content back to GitHub
func uploadGistChanges(gistID string, tempFilePath string) tea.Cmd {
	return func() tea.Msg {
		// Use gh CLI to edit the gist
		cmd := exec.Command("gh", "gist", "edit", gistID, "-a", tempFilePath)
		if err := cmd.Run(); err != nil {
			return errMsg{err: fmt.Errorf("failed to upload gist changes: %w", err)}
		}

		// Clean up temp file
		os.Remove(tempFilePath)

		// Reload gists to show updated timestamp
		return fetchGists()()
	}
}

// createGistFromFile creates a new gist from a temp file
func createGistFromFile(tempFilePath string, description string, public bool) tea.Cmd {
	return func() tea.Msg {
		// Prepare gh CLI command
		args := []string{"gist", "create", tempFilePath}

		if description != "" {
			args = append(args, "-d", description)
		}

		if public {
			args = append(args, "-p")
		}

		cmd := exec.Command("gh", args...)
		if err := cmd.Run(); err != nil {
			os.Remove(tempFilePath)
			return errMsg{err: fmt.Errorf("failed to create gist: %w", err)}
		}

		// Clean up temp file
		os.Remove(tempFilePath)

		// Reload gists to show new gist
		return fetchGists()()
	}
}
