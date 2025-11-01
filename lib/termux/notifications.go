package termux

import (
	"os/exec"
	"strings"
)

// NotifyOption is a functional option for configuring notifications.
type NotifyOption func(*notifyConfig)

// notifyConfig holds the configuration for a notification.
type notifyConfig struct {
	id       string
	ongoing  bool
	priority string
	icon     string
	vibrate  string
	sound    bool
	buttons  []button
}

// button represents an action button in a notification.
type button struct {
	text   string
	action string
}

// Notify displays an Android notification with the given title and content.
// Additional options can be provided to customize the notification behavior.
//
// If not running on Termux, this is a no-op.
//
// Example:
//
//	termux.Notify("Task Complete", "Build finished successfully")
//
//	termux.Notify("AI Worker", "Processing task #123",
//	    termux.WithID("worker-123"),
//	    termux.WithPriority("high"),
//	    termux.WithButton("View", "termux-open-url https://..."),
//	)
func Notify(title, content string, opts ...NotifyOption) error {
	if !IsTermux() {
		return nil
	}

	config := &notifyConfig{}
	for _, opt := range opts {
		opt(config)
	}

	args := []string{"termux-notification"}
	args = append(args, "--title", title)
	args = append(args, "--content", content)

	if config.id != "" {
		args = append(args, "--id", config.id)
	}

	if config.ongoing {
		args = append(args, "--ongoing")
	}

	if config.priority != "" {
		args = append(args, "--priority", config.priority)
	}

	if config.icon != "" {
		args = append(args, "--icon", config.icon)
	}

	if config.vibrate != "" {
		args = append(args, "--vibrate", config.vibrate)
	}

	if config.sound {
		args = append(args, "--sound")
	}

	// Add up to 3 buttons
	buttonLabels := []string{"--button1", "--button2", "--button3"}
	actionLabels := []string{"--button1-action", "--button2-action", "--button3-action"}
	for i, btn := range config.buttons {
		if i >= 3 {
			break // Termux supports max 3 buttons
		}
		args = append(args, buttonLabels[i], btn.text)
		args = append(args, actionLabels[i], btn.action)
	}

	cmd := exec.Command(args[0], args[1:]...)
	return cmd.Run()
}

// NotifyRemove removes a notification by its ID.
//
// If not running on Termux, this is a no-op.
//
// Example:
//
//	termux.NotifyRemove("worker-123")
func NotifyRemove(id string) error {
	if !IsTermux() {
		return nil
	}

	cmd := exec.Command("termux-notification-remove", id)
	return cmd.Run()
}

// WithID sets a unique identifier for the notification.
// This allows updating or removing the notification later.
//
// Example:
//
//	termux.Notify("Worker", "Processing...",
//	    termux.WithID("worker-123"))
//
//	// Later, update it:
//	termux.Notify("Worker", "Complete!",
//	    termux.WithID("worker-123"))
//
//	// Or remove it:
//	termux.NotifyRemove("worker-123")
func WithID(id string) NotifyOption {
	return func(c *notifyConfig) {
		c.id = id
	}
}

// WithOngoing makes the notification persistent (cannot be swiped away).
// Useful for background tasks or ongoing operations.
//
// Example:
//
//	termux.Notify("Sync Active", "Syncing files...",
//	    termux.WithOngoing())
func WithOngoing() NotifyOption {
	return func(c *notifyConfig) {
		c.ongoing = true
	}
}

// WithPriority sets the notification priority level.
// Valid values: "default", "high", "low", "max", "min"
//
// High priority notifications appear at the top and may make sound/vibrate.
// Low priority notifications appear at the bottom and are silent.
//
// Example:
//
//	termux.Notify("Error", "Build failed",
//	    termux.WithPriority("high"))
func WithPriority(priority string) NotifyOption {
	return func(c *notifyConfig) {
		c.priority = priority
	}
}

// WithIcon sets the notification icon.
// Common icons: "sync", "error", "info", "warning"
//
// Example:
//
//	termux.Notify("Syncing", "Downloading files...",
//	    termux.WithIcon("sync"))
func WithIcon(icon string) NotifyOption {
	return func(c *notifyConfig) {
		c.icon = icon
	}
}

// WithVibrate sets a custom vibration pattern in milliseconds.
// Pattern format: "duration1,pause1,duration2,pause2,..."
//
// Example:
//
//	termux.Notify("Alert", "Important message",
//	    termux.WithVibrate("100,50,100"))  // Two short bursts
func WithVibrate(pattern string) NotifyOption {
	return func(c *notifyConfig) {
		c.vibrate = pattern
	}
}

// WithSound enables the notification sound.
//
// Example:
//
//	termux.Notify("Complete", "Task finished",
//	    termux.WithSound())
func WithSound() NotifyOption {
	return func(c *notifyConfig) {
		c.sound = true
	}
}

// WithButton adds an action button to the notification.
// Up to 3 buttons can be added. Additional buttons will be ignored.
//
// The action parameter is a command to execute when the button is tapped.
// Common actions:
//   - "termux-open-url <url>" - Open a URL
//   - "termux-clipboard-set <text>" - Copy text to clipboard
//   - Any shell command
//
// Example:
//
//	termux.Notify("PR Ready", "Pull request #123 created",
//	    termux.WithButton("View", "termux-open-url https://github.com/..."),
//	    termux.WithButton("Copy URL", "bash -c 'echo https://... | termux-clipboard-set'"))
func WithButton(text, action string) NotifyOption {
	return func(c *notifyConfig) {
		c.buttons = append(c.buttons, button{text: text, action: action})
	}
}

// NotificationList returns a list of currently active notifications.
// The output is in JSON format.
//
// If not running on Termux, returns an empty string.
func NotificationList() (string, error) {
	if !IsTermux() {
		return "", nil
	}

	cmd := exec.Command("termux-notification-list")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
