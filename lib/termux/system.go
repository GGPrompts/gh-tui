package termux

import (
	"encoding/json"
	"os/exec"
)

// BatteryStatus represents the current battery state of the device.
type BatteryStatus struct {
	Health      string  `json:"health"`       // Battery health status (e.g., "GOOD")
	Percentage  int     `json:"percentage"`   // Battery percentage (0-100)
	Plugged     string  `json:"plugged"`      // Plugged state (e.g., "PLUGGED_AC", "UNPLUGGED")
	Status      string  `json:"status"`       // Charging status (e.g., "CHARGING", "DISCHARGING")
	Temperature float64 `json:"temperature"`  // Battery temperature in Celsius
	Current     int     `json:"current"`      // Battery current in microamperes
	Voltage     int     `json:"voltage"`      // Battery voltage in millivolts
}

// GetBatteryStatus retrieves the current battery status.
//
// If not running on Termux, returns a default BatteryStatus with 100% battery.
//
// Example:
//
//	battery, err := termux.GetBatteryStatus()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Battery: %d%% (%s)\n", battery.Percentage, battery.Status)
//
//	// Battery-aware logic
//	if battery.Percentage < 20 && battery.Status != "CHARGING" {
//	    fmt.Println("Low battery - skipping heavy task")
//	    return
//	}
func GetBatteryStatus() (*BatteryStatus, error) {
	if !IsTermux() {
		return &BatteryStatus{
			Health:     "GOOD",
			Percentage: 100,
			Plugged:    "UNPLUGGED",
			Status:     "FULL",
		}, nil
	}

	cmd := exec.Command("termux-battery-status")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var status BatteryStatus
	if err := json.Unmarshal(output, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

// IsCharging is a convenience method to check if the device is charging.
//
// If not running on Termux, returns false.
func IsCharging() (bool, error) {
	status, err := GetBatteryStatus()
	if err != nil {
		return false, err
	}
	return status.Status == "CHARGING", nil
}

// IsBatteryLow checks if the battery is below the specified percentage.
//
// If not running on Termux, returns false.
//
// Example:
//
//	low, err := termux.IsBatteryLow(20)
//	if low {
//	    termux.Toast("Battery low - conserving power")
//	}
func IsBatteryLow(threshold int) (bool, error) {
	status, err := GetBatteryStatus()
	if err != nil {
		return false, err
	}
	return status.Percentage < threshold, nil
}

// Location represents GPS coordinates and location metadata.
type Location struct {
	Latitude  float64 `json:"latitude"`   // Latitude in degrees
	Longitude float64 `json:"longitude"`  // Longitude in degrees
	Altitude  float64 `json:"altitude"`   // Altitude in meters above sea level
	Accuracy  float64 `json:"accuracy"`   // Horizontal accuracy in meters
	Bearing   float64 `json:"bearing"`    // Direction of travel in degrees (0-360)
	Speed     float64 `json:"speed"`      // Speed in meters per second
	Provider  string  `json:"provider"`   // Location provider (e.g., "gps", "network")
}

// GetLocation retrieves the current GPS location.
// This may take a few seconds to acquire a GPS fix.
//
// If not running on Termux, returns a default Location at (0, 0).
//
// Example:
//
//	loc, err := termux.GetLocation()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Location: %.4f, %.4f (accuracy: %.1fm)\n",
//	    loc.Latitude, loc.Longitude, loc.Accuracy)
func GetLocation() (*Location, error) {
	if !IsTermux() {
		return &Location{}, nil
	}

	cmd := exec.Command("termux-location")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var loc Location
	if err := json.Unmarshal(output, &loc); err != nil {
		return nil, err
	}

	return &loc, nil
}

// GetLocationWithProvider retrieves location using a specific provider.
//
// Providers:
//   - "gps" - GPS hardware (most accurate, outdoor only)
//   - "network" - Cell tower/WiFi triangulation (faster, less accurate)
//   - "passive" - Use last known location without requesting new fix
//
// If not running on Termux, returns a default Location at (0, 0).
//
// Example:
//
//	loc, err := termux.GetLocationWithProvider("network")  // Faster than GPS
func GetLocationWithProvider(provider string) (*Location, error) {
	if !IsTermux() {
		return &Location{}, nil
	}

	cmd := exec.Command("termux-location", "-p", provider)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var loc Location
	if err := json.Unmarshal(output, &loc); err != nil {
		return nil, err
	}

	return &loc, nil
}

// SensorData represents data from a device sensor.
type SensorData struct {
	Sensor string                 `json:"sensor"` // Sensor name
	Values map[string]interface{} `json:"values"` // Sensor readings
}

// GetSensor retrieves data from a specific sensor.
//
// Common sensors:
//   - "accelerometer" - Motion detection (x, y, z acceleration)
//   - "light" - Ambient light level (lux)
//   - "proximity" - Distance to nearby objects (cm)
//   - "gyroscope" - Rotation rate (x, y, z)
//   - "magnetic_field" - Magnetic field strength (compass)
//   - "pressure" - Barometric pressure (hPa)
//   - "temperature" - Device temperature (°C)
//   - "humidity" - Relative humidity (%)
//
// If not running on Termux, returns an empty SensorData.
//
// Example:
//
//	sensor, err := termux.GetSensor("light")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Check light level for automatic brightness adjustment
func GetSensor(sensorName string) (*SensorData, error) {
	if !IsTermux() {
		return &SensorData{Sensor: sensorName, Values: map[string]interface{}{}}, nil
	}

	cmd := exec.Command("termux-sensor", "-s", sensorName, "-n", "1")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var data SensorData
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// ListSensors returns a list of available sensors on the device.
//
// If not running on Termux, returns an empty string.
//
// Example:
//
//	sensors, err := termux.ListSensors()
//	fmt.Println("Available sensors:", sensors)
func ListSensors() (string, error) {
	if !IsTermux() {
		return "", nil
	}

	cmd := exec.Command("termux-sensor", "-l")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// WiFiConnectionInfo represents the current WiFi connection details.
type WiFiConnectionInfo struct {
	SSID          string `json:"ssid"`            // Network name
	BSSID         string `json:"bssid"`           // Access point MAC address
	IP            string `json:"ip"`              // Device IP address
	MAC           string `json:"mac"`             // Device MAC address
	RSSI          int    `json:"rssi"`            // Signal strength (dBm)
	LinkSpeedMbps int    `json:"link_speed_mbps"` // Connection speed (Mbps)
	FrequencyMhz  int    `json:"frequency_mhz"`   // WiFi frequency (MHz)
}

// GetWiFiConnectionInfo retrieves information about the current WiFi connection.
//
// If not running on Termux or not connected to WiFi, returns an empty struct.
//
// Example:
//
//	wifi, err := termux.GetWiFiConnectionInfo()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if wifi.SSID == "HomeNetwork" {
//	    // Safe to run automation at home
//	}
func GetWiFiConnectionInfo() (*WiFiConnectionInfo, error) {
	if !IsTermux() {
		return &WiFiConnectionInfo{}, nil
	}

	cmd := exec.Command("termux-wifi-connectioninfo")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var info WiFiConnectionInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

// WiFiScanResult represents a WiFi network found during scanning.
type WiFiScanResult struct {
	SSID         string `json:"ssid"`          // Network name
	BSSID        string `json:"bssid"`         // Access point MAC address
	RSSI         int    `json:"rssi"`          // Signal strength (dBm)
	FrequencyMhz int    `json:"frequency_mhz"` // WiFi frequency (MHz)
}

// ScanWiFi scans for available WiFi networks.
//
// If not running on Termux, returns an empty slice.
//
// Example:
//
//	networks, err := termux.ScanWiFi()
//	for _, net := range networks {
//	    fmt.Printf("%s: %d dBm\n", net.SSID, net.RSSI)
//	}
func ScanWiFi() ([]WiFiScanResult, error) {
	if !IsTermux() {
		return []WiFiScanResult{}, nil
	}

	cmd := exec.Command("termux-wifi-scaninfo")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var results []WiFiScanResult
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// SetWiFiEnabled enables or disables WiFi.
//
// If not running on Termux, this is a no-op.
//
// Example:
//
//	termux.SetWiFiEnabled(true)   // Enable WiFi
//	termux.SetWiFiEnabled(false)  // Disable WiFi
func SetWiFiEnabled(enabled bool) error {
	if !IsTermux() {
		return nil
	}

	value := "false"
	if enabled {
		value = "true"
	}

	cmd := exec.Command("termux-wifi-enable", value)
	return cmd.Run()
}

// ClipboardSet copies text to the Android clipboard.
// This uses the heredoc pattern from TFE's clipboard implementation
// to avoid exit status 2 errors with direct StdinPipe.
//
// If not running on Termux, this is a no-op.
//
// Example:
//
//	termux.ClipboardSet("https://github.com/user/repo/pull/123")
func ClipboardSet(text string) error {
	if !IsTermux() {
		return nil
	}

	// Use heredoc pattern from TFE editor.go:115
	cmd := exec.Command("bash", "-c", "termux-clipboard-set <<'CLIPBOARD_EOF'\n"+text+"\nCLIPBOARD_EOF")
	return cmd.Run()
}

// ClipboardGet retrieves text from the Android clipboard.
//
// If not running on Termux, returns an empty string.
//
// Example:
//
//	text, err := termux.ClipboardGet()
//	fmt.Println("Clipboard:", text)
func ClipboardGet() (string, error) {
	if !IsTermux() {
		return "", nil
	}

	cmd := exec.Command("termux-clipboard-get")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// WakeLock acquires a wake lock to prevent the device from sleeping.
// This is essential for long-running background tasks.
//
// IMPORTANT: Always release the wake lock when done using WakeUnlock()
// or defer it to ensure it's released even on error.
//
// If not running on Termux, this is a no-op.
//
// Example:
//
//	termux.WakeLock()
//	defer termux.WakeUnlock()
//	// Do long-running work...
func WakeLock() error {
	if !IsTermux() {
		return nil
	}

	cmd := exec.Command("termux-wake-lock")
	return cmd.Run()
}

// WakeUnlock releases the wake lock, allowing the device to sleep normally.
//
// If not running on Termux, this is a no-op.
func WakeUnlock() error {
	if !IsTermux() {
		return nil
	}

	cmd := exec.Command("termux-wake-unlock")
	return cmd.Run()
}
