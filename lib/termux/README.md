# Termux API Library for Go

A comprehensive, production-ready Go wrapper for Termux API commands. Enables TUI applications to integrate seamlessly with Android device features when running in Termux.

## Features

- **Zero Dependencies** - Only uses Go standard library
- **Graceful Degradation** - All functions are no-ops when not on Termux (safe for cross-platform apps)
- **Type-Safe** - Proper structs for all JSON responses
- **Production Ready** - Clean error handling, comprehensive documentation
- **Complete Coverage** - All major Termux API features included

## Installation

```bash
# Install Termux API package on your Android device
pkg install termux-api

# Also install Termux:API app from F-Droid:
# https://f-droid.org/packages/com.termux.api/

# In your Go project
go get github.com/yourname/TUITemplate/lib/termux
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/yourname/TUITemplate/lib/termux"
)

func main() {
    // Detect Termux environment
    if termux.IsTermux() {
        termux.Toast("Running on Termux!")
        termux.Vibrate(50)
    }

    // Show notification with action buttons
    termux.Notify(
        "Task Complete",
        "Your build finished successfully",
        termux.WithButton("View", "termux-open-url https://..."),
        termux.WithSound(),
    )

    // Check battery before heavy work
    battery, _ := termux.GetBatteryStatus()
    if battery.Percentage < 20 && battery.Status != "CHARGING" {
        termux.Toast("Low battery - skipping task")
        return
    }

    // Get user confirmation
    confirmed, _ := termux.ConfirmDialog("Continue?", "Run heavy task?")
    if confirmed {
        // Do work...
    }
}
```

## API Documentation

### Core Functions

#### Detection

```go
// Check if running on Termux (cached for performance)
if termux.IsTermux() {
    // Termux-specific code
}
```

#### Haptic Feedback

```go
// Quick vibration (50ms) for button press
termux.Vibrate(50)

// Medium vibration for success
termux.Vibrate(100)

// Long vibration for error
termux.Vibrate(500)

// Force vibration (even in silent mode)
termux.VibrateForce(100)
```

#### Toast Messages

```go
// Short toast (~2 seconds)
termux.Toast("File saved")

// Long toast (~4 seconds)
termux.ToastLong("Processing complete - 15 files updated")

// Very short toast
termux.ToastShort("Done")
```

### Notifications

Rich Android notifications with action buttons, priority levels, and custom vibration patterns.

#### Basic Notification

```go
termux.Notify("Title", "Content")
```

#### Advanced Notification

```go
termux.Notify(
    "Pull Request Ready",
    "PR #123 has been created",
    termux.WithID("pr-123"),           // Unique ID for updates/removal
    termux.WithPriority("high"),       // Priority: default, high, low, max, min
    termux.WithIcon("sync"),           // Icon: sync, error, info, warning
    termux.WithSound(),                // Play notification sound
    termux.WithVibrate("100,50,100"),  // Custom vibration pattern
    termux.WithButton("View", "termux-open-url https://github.com/..."),
    termux.WithButton("Copy", "bash -c 'echo URL | termux-clipboard-set'"),
)
```

#### Ongoing Notification (Persistent)

```go
// Show ongoing notification (can't swipe away)
termux.Notify(
    "Worker Active",
    "Processing tasks...",
    termux.WithID("worker"),
    termux.WithOngoing(),
)

// Update it
termux.Notify(
    "Worker Active",
    "5 of 10 tasks complete",
    termux.WithID("worker"),
    termux.WithOngoing(),
)

// Remove when done
termux.NotifyRemove("worker")
```

### Voice & Speech

#### Speech to Text

```go
termux.Toast("Listening...")
text, err := termux.SpeechToText()
if err != nil {
    log.Fatal(err)
}
fmt.Println("You said:", text)
```

#### Text to Speech

```go
// Basic TTS
termux.Speak("Task complete")

// Advanced TTS with options
termux.SpeakWithOptions(
    "Your AI worker has finished",
    "",        // engine (empty = default)
    "en-US",   // language
    1.0,       // pitch
    1.0,       // rate
    "",        // stream (empty = default)
)
```

### Dialogs

Native Android dialogs for user input.

#### Confirmation Dialog

```go
confirmed, err := termux.ConfirmDialog(
    "Approve PR?",
    "Merge pull request #123?",
)
if confirmed {
    // User clicked "yes"
}
```

#### Text Input Dialog

```go
message, err := termux.TextDialog(
    "Commit Message",
    "Enter commit message:",
)
```

#### Password Dialog (Hidden Input)

```go
apiKey, err := termux.PasswordDialog(
    "API Key",
    "Enter your API key:",
)
```

#### Radio Dialog (Single Choice)

```go
action, err := termux.RadioDialog(
    "Choose Action",
    "Approve,Reject,Review,Cancel",
)
// Returns one of: "Approve", "Reject", "Review", or "Cancel"
```

#### Checkbox Dialog (Multiple Choice)

```go
options, err := termux.CheckboxDialog(
    "Select Options",
    "Run Tests,Build,Deploy,Notify",
)
// Returns slice: ["Run Tests", "Build", "Deploy"]
```

#### Spinner Dialog (Dropdown)

```go
model, err := termux.SpinnerDialog(
    "Choose Model",
    "sonnet,opus,haiku",
)
```

#### Date & Time Dialogs

```go
// Date picker
date, err := termux.DateDialog("Select Date", "2025-10-30")
// Returns: "2025-10-30"

// Time picker
time, err := termux.TimeDialog("Select Time", "14:30")
// Returns: "14:30"
```

#### Counter Dialog

```go
count, err := termux.CounterDialog("How many tasks?", 1, 10)
// Returns string representation of count
```

### Battery & Power

#### Battery Status

```go
battery, err := termux.GetBatteryStatus()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Battery: %d%%\n", battery.Percentage)
fmt.Printf("Status: %s\n", battery.Status)      // CHARGING, DISCHARGING, FULL
fmt.Printf("Plugged: %s\n", battery.Plugged)    // PLUGGED_AC, UNPLUGGED
fmt.Printf("Health: %s\n", battery.Health)      // GOOD, OVERHEAT, etc.
fmt.Printf("Temp: %.1f°C\n", battery.Temperature)
```

#### Battery-Aware Logic

```go
// Check if charging
charging, err := termux.IsCharging()
if charging {
    // Safe to run heavy task
}

// Check battery level
low, err := termux.IsBatteryLow(20)  // threshold: 20%
if low {
    termux.Toast("Battery low - skipping task")
    return
}

// Complete example
battery, _ := termux.GetBatteryStatus()
if battery.Percentage < 20 && battery.Status != "CHARGING" {
    termux.Notify(
        "Task Skipped",
        "Battery too low - connect charger",
        termux.WithPriority("high"),
    )
    return
}
```

#### Wake Lock (Prevent Sleep)

```go
// Acquire wake lock for long-running task
termux.WakeLock()
defer termux.WakeUnlock()  // Always release!

// Do long-running work...
```

### Location & GPS

#### Get Location

```go
loc, err := termux.GetLocation()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Lat: %.4f\n", loc.Latitude)
fmt.Printf("Lon: %.4f\n", loc.Longitude)
fmt.Printf("Alt: %.1fm\n", loc.Altitude)
fmt.Printf("Accuracy: ±%.1fm\n", loc.Accuracy)
fmt.Printf("Speed: %.1f m/s\n", loc.Speed)
fmt.Printf("Provider: %s\n", loc.Provider)  // gps, network
```

#### Location with Specific Provider

```go
// GPS (accurate, outdoor only, slower)
loc, err := termux.GetLocationWithProvider("gps")

// Network (cell tower/WiFi, faster, less accurate)
loc, err := termux.GetLocationWithProvider("network")

// Passive (last known location, instant)
loc, err := termux.GetLocationWithProvider("passive")
```

### WiFi & Network

#### WiFi Connection Info

```go
wifi, err := termux.GetWiFiConnectionInfo()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("SSID: %s\n", wifi.SSID)
fmt.Printf("IP: %s\n", wifi.IP)
fmt.Printf("Signal: %d dBm\n", wifi.RSSI)
fmt.Printf("Speed: %d Mbps\n", wifi.LinkSpeedMbps)

// Network-aware automation
if wifi.SSID == "HomeNetwork" {
    // Safe to run at home
}
```

#### WiFi Scanning

```go
networks, err := termux.ScanWiFi()
for _, net := range networks {
    fmt.Printf("%s: %d dBm\n", net.SSID, net.RSSI)
}
```

#### WiFi Control

```go
// Enable WiFi
termux.SetWiFiEnabled(true)

// Disable WiFi
termux.SetWiFiEnabled(false)
```

### Sensors

#### Read Sensor Data

```go
// Light sensor (for automatic brightness)
light, err := termux.GetSensor("light")
fmt.Printf("Light: %v\n", light.Values)

// Proximity sensor
proximity, err := termux.GetSensor("proximity")

// Accelerometer (motion detection)
accel, err := termux.GetSensor("accelerometer")

// List all available sensors
sensors, err := termux.ListSensors()
fmt.Println(sensors)
```

Available sensors:
- `accelerometer` - Motion detection (x, y, z)
- `light` - Ambient light level (lux)
- `proximity` - Distance to nearby objects (cm)
- `gyroscope` - Rotation rate
- `magnetic_field` - Compass
- `pressure` - Barometric pressure (hPa)
- `temperature` - Device temperature (°C)
- `humidity` - Relative humidity (%)

### Clipboard

```go
// Copy to clipboard
url := "https://github.com/user/repo/pull/123"
termux.ClipboardSet(url)
termux.Toast("URL copied")

// Read from clipboard
text, err := termux.ClipboardGet()
fmt.Println("Clipboard:", text)
```

## Usage Examples

### Haptic Feedback in TUI List

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/yourname/TUITemplate/lib/termux"
)

type model struct {
    items  []string
    cursor int
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up":
            if m.cursor > 0 {
                m.cursor--
                termux.Vibrate(30)  // Quick haptic
            }
        case "down":
            if m.cursor < len(m.items)-1 {
                m.cursor++
                termux.Vibrate(30)
            }
        case "enter":
            termux.Vibrate(50)
            termux.Toast("Selected: " + m.items[m.cursor])
            // Process selection...
        }
    }
    return m, nil
}
```

### Progress Notifications

```go
func ProcessTasks(tasks []string) {
    for i, task := range tasks {
        // Update progress
        termux.Notify(
            "Processing Tasks",
            fmt.Sprintf("%d/%d: %s", i+1, len(tasks), task),
            termux.WithID("progress"),
            termux.WithOngoing(),
        )

        // Do work...
        processTask(task)
    }

    // Final notification
    termux.Notify(
        "Complete",
        fmt.Sprintf("All %d tasks finished", len(tasks)),
        termux.WithID("progress"),
        termux.WithSound(),
        termux.WithVibrate("100,50,100"),
    )

    termux.Speak("All tasks complete")
}
```

### Voice-Controlled TUI

```go
func VoiceCommand() {
    termux.Toast("Listening...")

    text, err := termux.SpeechToText()
    if err != nil {
        termux.Speak("Voice input failed")
        return
    }

    termux.Vibrate(50)

    switch {
    case strings.Contains(text, "sync"):
        termux.Speak("Syncing projects")
        syncProjects()
    case strings.Contains(text, "status"):
        termux.Speak("All systems operational")
        showStatus()
    case strings.Contains(text, "exit"):
        termux.Speak("Goodbye")
        os.Exit(0)
    default:
        termux.Speak("Unknown command")
    }
}
```

### Battery-Aware Background Worker

```go
func BackgroundWorker() {
    // Acquire wake lock
    termux.WakeLock()
    defer termux.WakeUnlock()

    // Check battery
    battery, _ := termux.GetBatteryStatus()
    if battery.Percentage < 20 && battery.Status != "CHARGING" {
        termux.Notify(
            "Worker Skipped",
            fmt.Sprintf("Battery too low: %d%%", battery.Percentage),
            termux.WithPriority("low"),
        )
        return
    }

    // Check network
    wifi, _ := termux.GetWiFiConnectionInfo()
    if wifi.SSID != "HomeNetwork" {
        termux.Toast("Not on trusted network")
        return
    }

    // Show ongoing notification
    termux.Notify(
        "Worker Active",
        "Processing tasks...",
        termux.WithID("worker"),
        termux.WithOngoing(),
    )

    // Do work
    tasks := fetchTasks()
    successCount := 0

    for _, task := range tasks {
        termux.Toast(fmt.Sprintf("Processing: %s", task.Title))

        if processTask(task) {
            successCount++
            termux.Vibrate(50)
        } else {
            termux.Vibrate(200)
        }
    }

    // Final notification
    termux.Notify(
        "Worker Complete",
        fmt.Sprintf("%d tasks completed", successCount),
        termux.WithID("worker-done"),
        termux.WithPriority("high"),
        termux.WithSound(),
        termux.WithButton("View Results", "termux-open-url https://..."),
    )

    termux.Speak(fmt.Sprintf("Completed %d tasks", successCount))
    termux.NotifyRemove("worker")
}
```

### Interactive Task Approval

```go
func ReviewPullRequest(prNumber int) {
    // Show PR info
    pr := getPRInfo(prNumber)

    // Ask for action
    action, err := termux.RadioDialog(
        "Pull Request #" + strconv.Itoa(prNumber),
        "Approve,Request Changes,Comment,Close",
    )
    if err != nil {
        log.Fatal(err)
    }

    switch action {
    case "Approve":
        confirmed, _ := termux.ConfirmDialog(
            "Approve PR?",
            "This will merge the pull request",
        )
        if confirmed {
            termux.Vibrate(100)
            approvePR(prNumber)
            termux.Notify(
                "PR Approved",
                fmt.Sprintf("PR #%d merged", prNumber),
                termux.WithSound(),
            )
        }

    case "Request Changes":
        comment, _ := termux.TextDialog(
            "Changes Requested",
            "Enter your feedback:",
        )
        if comment != "" {
            requestChanges(prNumber, comment)
            termux.Toast("Feedback sent")
        }

    case "Comment":
        comment, _ := termux.TextDialog("Comment", "Enter comment:")
        if comment != "" {
            addComment(prNumber, comment)
            termux.Toast("Comment added")
        }

    case "Close":
        closePR(prNumber)
        termux.Toast("PR closed")
    }
}
```

## Best Practices

### 1. Always Check IsTermux() for Optional Features

```go
// Good: Graceful degradation
if termux.IsTermux() {
    termux.Vibrate(50)
    termux.Toast("Action complete")
}

// Also good: All functions are safe to call
termux.Vibrate(50)  // No-op if not on Termux
```

### 2. Release Wake Locks with defer

```go
termux.WakeLock()
defer termux.WakeUnlock()  // Ensures release even on error

// Do work...
```

### 3. Check Battery Before Heavy Operations

```go
battery, _ := termux.GetBatteryStatus()
if battery.Percentage < 20 && battery.Status != "CHARGING" {
    // Skip or defer heavy task
    return
}
```

### 4. Use Notification IDs for Updates

```go
// Initial
termux.Notify("Syncing", "Starting...", termux.WithID("sync"), termux.WithOngoing())

// Update
termux.Notify("Syncing", "50% complete", termux.WithID("sync"), termux.WithOngoing())

// Remove
termux.NotifyRemove("sync")
```

### 5. Provide User Feedback

```go
// Quick feedback for interactions
termux.Vibrate(30)      // Haptic
termux.Toast("Saved")   // Visual

// Important events
termux.Notify("Complete", "Task finished", termux.WithSound())
termux.Speak("Task complete")  // Audio
```

## Integration with TUITemplate

This library is designed to work seamlessly with TUITemplate projects:

```go
// In your TUI update function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            termux.Vibrate(30)  // Haptic on panel switch
            m.focusedPanel = switchPanel(m.focusedPanel)

        case "enter":
            termux.Vibrate(50)
            selected := m.getSelectedItem()
            termux.Toast("Opening: " + selected.Name)
            // Process selection...
        }
    }
    return m, nil
}
```

## Future Projects Using This Library

### tmuxplexer
- Haptic feedback on panel switch
- Notifications when Claude Code finishes tasks
- Voice commands to switch sessions

### tkan (Kanban TUI)
- Vibrate on card drag/drop
- Battery-aware GitHub API sync
- Voice input for creating cards
- Location-based task filtering

### TFE (File Manager)
- Clipboard integration for file paths
- Voice commands for navigation
- Battery status in status bar

## Requirements

- Termux app (from F-Droid or Play Store)
- Termux:API app (from F-Droid)
- `pkg install termux-api` in Termux

## Platform Support

- **Android (Termux)**: Full functionality
- **Linux/macOS/Windows**: All functions gracefully degrade to no-ops (safe for cross-platform apps)

## Testing

See `examples_test.go` for comprehensive usage examples. Run examples with:

```bash
go test -v ./lib/termux
```

## License

MIT License - See LICENSE file for details

## Contributing

Contributions welcome! Please ensure:
- All functions gracefully degrade when not on Termux
- Proper error handling throughout
- Documentation with examples
- Type-safe JSON parsing

## Credits

Created for the TUITemplate project. Inspired by:
- TFE clipboard integration pattern
- Termux API Reference Guide
- LazyGit's weight-based layout system
