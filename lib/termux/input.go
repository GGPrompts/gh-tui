package termux

import (
	"encoding/json"
	"os/exec"
	"strings"
)

// DialogResult represents the result from a dialog interaction.
type DialogResult struct {
	Code   int      `json:"code"`
	Text   string   `json:"text"`
	Values []string `json:"values,omitempty"` // For checkbox dialogs
}

// SpeechToText converts speech to text using Google's speech recognition.
// The function blocks until speech is detected and processed.
//
// If not running on Termux, returns an empty string.
//
// Example:
//
//	text, err := termux.SpeechToText()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("You said:", text)
func SpeechToText() (string, error) {
	if !IsTermux() {
		return "", nil
	}

	cmd := exec.Command("termux-speech-to-text")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// Speak converts text to speech using the device's TTS engine.
// The function returns immediately (speech plays asynchronously).
//
// If not running on Termux, this is a no-op.
//
// Example:
//
//	termux.Speak("Task complete")
func Speak(text string) error {
	if !IsTermux() {
		return nil
	}

	cmd := exec.Command("termux-tts-speak", text)
	return cmd.Run()
}

// SpeakWithOptions speaks text with advanced TTS options.
//
// Parameters:
//   - text: The text to speak
//   - engine: TTS engine (e.g., "com.google.android.tts")
//   - language: Language code (e.g., "en-US", "es-ES")
//   - pitch: Pitch level (0.0 - 2.0, default 1.0)
//   - rate: Speech rate (0.0 - 2.0, default 1.0)
//   - stream: Audio stream type (e.g., "STREAM_NOTIFICATION")
//
// If not running on Termux, this is a no-op.
//
// Example:
//
//	termux.SpeakWithOptions("Hello world", "", "en-US", 1.2, 0.9, "")
func SpeakWithOptions(text, engine, language string, pitch, rate float64, stream string) error {
	if !IsTermux() {
		return nil
	}

	args := []string{"termux-tts-speak"}

	if engine != "" {
		args = append(args, "-e", engine)
	}
	if language != "" {
		args = append(args, "-l", language)
	}
	if pitch != 0.0 {
		args = append(args, "-p", formatFloat(pitch))
	}
	if rate != 0.0 {
		args = append(args, "-r", formatFloat(rate))
	}
	if stream != "" {
		args = append(args, "-s", stream)
	}

	args = append(args, text)

	cmd := exec.Command(args[0], args[1:]...)
	return cmd.Run()
}

// Dialog shows a native Android dialog and returns the user's input.
//
// Parameters:
//   - dialogType: Type of dialog ("confirm", "text", "radio", "checkbox", "spinner", "date", "time", "counter")
//   - title: Dialog title
//   - hint: Hint/prompt text
//
// Additional parameters vary by dialog type. Use the type-specific helper functions instead.
//
// If not running on Termux, returns an empty DialogResult.
func Dialog(dialogType, title, hint string) (*DialogResult, error) {
	if !IsTermux() {
		return &DialogResult{}, nil
	}

	args := []string{"termux-dialog", dialogType, "-t", title}
	if hint != "" {
		args = append(args, "-i", hint)
	}

	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var result DialogResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ConfirmDialog shows a yes/no confirmation dialog.
//
// If not running on Termux, returns false.
//
// Example:
//
//	confirmed, err := termux.ConfirmDialog("Approve PR?", "Merge pull request #123?")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if confirmed {
//	    // User clicked "yes"
//	}
func ConfirmDialog(title, message string) (bool, error) {
	if !IsTermux() {
		return false, nil
	}

	result, err := Dialog("confirm", title, message)
	if err != nil {
		return false, err
	}

	return result.Text == "yes", nil
}

// TextDialog shows a text input dialog and returns the entered text.
//
// If not running on Termux, returns an empty string.
//
// Example:
//
//	message, err := termux.TextDialog("Commit Message", "Enter commit message:")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Message:", message)
func TextDialog(title, hint string) (string, error) {
	if !IsTermux() {
		return "", nil
	}

	result, err := Dialog("text", title, hint)
	if err != nil {
		return "", err
	}

	return result.Text, nil
}

// PasswordDialog shows a password input dialog (hidden text entry).
//
// If not running on Termux, returns an empty string.
//
// Example:
//
//	apiKey, err := termux.PasswordDialog("API Key", "Enter your API key:")
func PasswordDialog(title, hint string) (string, error) {
	if !IsTermux() {
		return "", nil
	}

	args := []string{"termux-dialog", "text", "-t", title, "-i", hint, "-p"}
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var result DialogResult
	if err := json.Unmarshal(output, &result); err != nil {
		return "", err
	}

	return result.Text, nil
}

// RadioDialog shows a radio button dialog (single choice).
//
// Parameters:
//   - title: Dialog title
//   - values: Comma-separated list of options
//
// If not running on Termux, returns an empty string.
//
// Example:
//
//	action, err := termux.RadioDialog("Choose Action", "Approve,Reject,Review,Cancel")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Selected:", action)
func RadioDialog(title string, values string) (string, error) {
	if !IsTermux() {
		return "", nil
	}

	args := []string{"termux-dialog", "radio", "-t", title, "-v", values}
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var result DialogResult
	if err := json.Unmarshal(output, &result); err != nil {
		return "", err
	}

	return result.Text, nil
}

// CheckboxDialog shows a checkbox dialog (multiple choice).
//
// Parameters:
//   - title: Dialog title
//   - values: Comma-separated list of options
//
// Returns a slice of selected values.
//
// If not running on Termux, returns an empty slice.
//
// Example:
//
//	options, err := termux.CheckboxDialog("Select Options", "Run Tests,Build,Deploy,Notify")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, opt := range options {
//	    fmt.Println("Selected:", opt)
//	}
func CheckboxDialog(title string, values string) ([]string, error) {
	if !IsTermux() {
		return []string{}, nil
	}

	args := []string{"termux-dialog", "checkbox", "-t", title, "-v", values}
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var result DialogResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	return result.Values, nil
}

// SpinnerDialog shows a dropdown spinner dialog.
//
// Parameters:
//   - title: Dialog title
//   - values: Comma-separated list of options
//
// If not running on Termux, returns an empty string.
//
// Example:
//
//	model, err := termux.SpinnerDialog("Choose Model", "sonnet,opus,haiku")
func SpinnerDialog(title string, values string) (string, error) {
	if !IsTermux() {
		return "", nil
	}

	args := []string{"termux-dialog", "spinner", "-t", title, "-v", values}
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var result DialogResult
	if err := json.Unmarshal(output, &result); err != nil {
		return "", err
	}

	return result.Text, nil
}

// DateDialog shows a date picker dialog.
//
// Parameters:
//   - title: Dialog title
//   - defaultDate: Default date in YYYY-MM-DD format (empty for today)
//
// Returns the selected date in YYYY-MM-DD format.
//
// If not running on Termux, returns an empty string.
//
// Example:
//
//	date, err := termux.DateDialog("Select Date", "2025-10-30")
func DateDialog(title, defaultDate string) (string, error) {
	if !IsTermux() {
		return "", nil
	}

	args := []string{"termux-dialog", "date", "-t", title}
	if defaultDate != "" {
		args = append(args, "-d", defaultDate)
	}

	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var result DialogResult
	if err := json.Unmarshal(output, &result); err != nil {
		return "", err
	}

	return result.Text, nil
}

// TimeDialog shows a time picker dialog.
//
// Parameters:
//   - title: Dialog title
//   - defaultTime: Default time in HH:MM format (empty for now)
//
// Returns the selected time in HH:MM format.
//
// If not running on Termux, returns an empty string.
//
// Example:
//
//	time, err := termux.TimeDialog("Select Time", "14:30")
func TimeDialog(title, defaultTime string) (string, error) {
	if !IsTermux() {
		return "", nil
	}

	args := []string{"termux-dialog", "time", "-t", title}
	if defaultTime != "" {
		args = append(args, "-d", defaultTime)
	}

	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var result DialogResult
	if err := json.Unmarshal(output, &result); err != nil {
		return "", err
	}

	return result.Text, nil
}

// CounterDialog shows a counter dialog with increment/decrement buttons.
//
// Parameters:
//   - title: Dialog title
//   - min: Minimum value
//   - max: Maximum value
//
// Returns the selected count as a string.
//
// If not running on Termux, returns "0".
//
// Example:
//
//	count, err := termux.CounterDialog("How many tasks?", 1, 10)
func CounterDialog(title string, min, max int) (string, error) {
	if !IsTermux() {
		return "0", nil
	}

	rangeStr := formatInt(min) + "," + formatInt(max)
	args := []string{"termux-dialog", "counter", "-t", title, "-r", rangeStr}

	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var result DialogResult
	if err := json.Unmarshal(output, &result); err != nil {
		return "", err
	}

	return result.Text, nil
}

// formatFloat converts a float64 to a string with reasonable precision.
func formatFloat(f float64) string {
	// Simple implementation for common values
	if f == 0.0 {
		return "0.0"
	}
	if f == 1.0 {
		return "1.0"
	}
	if f == 2.0 {
		return "2.0"
	}
	// For other values, approximate with 1 decimal
	intPart := int(f)
	fracPart := int((f - float64(intPart)) * 10)
	return formatInt(intPart) + "." + formatInt(fracPart)
}
