// Package termux provides a Go wrapper for Termux API commands.
// It enables TUI applications to integrate with Android device features
// when running in the Termux environment.
//
// All functions gracefully degrade to no-ops when not running on Termux,
// making it safe to use in cross-platform applications.
package termux

import (
	"os/exec"
	"sync"
)

var (
	// isTermuxEnv caches the result of Termux environment detection
	isTermuxEnv     *bool
	isTermuxEnvOnce sync.Once
)

// IsTermux checks if the application is running in a Termux environment.
// It caches the result for subsequent calls.
//
// Detection is done by checking if the termux-vibrate command is available,
// which is a lightweight command that should exist in all Termux installations
// with the Termux:API package installed.
func IsTermux() bool {
	isTermuxEnvOnce.Do(func() {
		result := commandAvailable("termux-vibrate")
		isTermuxEnv = &result
	})
	return *isTermuxEnv
}

// commandAvailable checks if a command is available in the system PATH.
// This is similar to the editorAvailable function from TFE's editor.go.
func commandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// Vibrate triggers phone vibration for the specified duration in milliseconds.
// This provides haptic feedback for user interactions.
//
// Common use cases:
//   - Short vibration (50ms) for button presses
//   - Medium vibration (100-200ms) for success actions
//   - Long vibration (500ms) for errors or important events
//
// If not running on Termux, this is a no-op.
//
// Example:
//
//	termux.Vibrate(50)  // Quick haptic feedback
func Vibrate(durationMs int) error {
	if !IsTermux() {
		return nil
	}

	cmd := exec.Command("termux-vibrate", "-d", formatInt(durationMs))
	return cmd.Run()
}

// VibrateForce triggers forced vibration that works even when the device
// is in silent mode. Use sparingly for critical alerts only.
//
// If not running on Termux, this is a no-op.
func VibrateForce(durationMs int) error {
	if !IsTermux() {
		return nil
	}

	cmd := exec.Command("termux-vibrate", "-d", formatInt(durationMs), "-f")
	return cmd.Run()
}

// Toast displays a short popup message (approximately 2 seconds).
// Toasts are non-intrusive and useful for quick status updates.
//
// If not running on Termux, this is a no-op.
//
// Example:
//
//	termux.Toast("File saved successfully")
func Toast(message string) error {
	if !IsTermux() {
		return nil
	}

	cmd := exec.Command("termux-toast", message)
	return cmd.Run()
}

// ToastLong displays a longer popup message (approximately 4 seconds).
// Use for messages that require slightly more reading time.
//
// If not running on Termux, this is a no-op.
//
// Example:
//
//	termux.ToastLong("Processing complete - 15 files updated")
func ToastLong(message string) error {
	if !IsTermux() {
		return nil
	}

	cmd := exec.Command("termux-toast", "-l", message)
	return cmd.Run()
}

// ToastShort displays a very brief popup message.
//
// If not running on Termux, this is a no-op.
func ToastShort(message string) error {
	if !IsTermux() {
		return nil
	}

	cmd := exec.Command("termux-toast", "-s", message)
	return cmd.Run()
}

// formatInt converts an integer to a string.
// Helper function to avoid importing strconv just for this.
func formatInt(n int) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	// Convert to string by building digits in reverse
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}

	// Reverse the digits
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}
