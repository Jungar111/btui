package bluetooth

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// DiscoveredDevice represents a device found during scanning
type DiscoveredDevice struct {
	BluetoothDevice
	RSSI      int
	Timestamp time.Time
}

// DiscoveryScanner manages the bluetoothctl scanning process
type DiscoveryScanner struct {
	cmd            *exec.Cmd
	ctx            context.Context
	cancel         context.CancelFunc
	discoveredDevs map[string]DiscoveredDevice
	mutex          sync.RWMutex
	isScanning     bool
}

// DiscoveryUpdateMsg contains discovered devices
type DiscoveryUpdateMsg struct {
	Devices []DiscoveredDevice
	Err     error
}

// NewDiscoveryScanner creates a new discovery scanner
func NewDiscoveryScanner() *DiscoveryScanner {
	return &DiscoveryScanner{
		discoveredDevs: make(map[string]DiscoveredDevice),
	}
}

// StartDiscovery begins the scanning process
func (ds *DiscoveryScanner) StartDiscovery() error {
	if ds.isScanning {
		return fmt.Errorf("discovery already running")
	}

	ds.ctx, ds.cancel = context.WithCancel(context.Background())
	ds.cmd = exec.CommandContext(ds.ctx, "bluetoothctl")

	stdin, err := ds.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	stdout, err := ds.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := ds.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start bluetoothctl: %w", err)
	}

	ds.isScanning = true

	// Start scanning
	fmt.Fprintln(stdin, "scan on")

	// Parse output in a goroutine
	go func() {
		defer stdin.Close()
		scanner := bufio.NewScanner(stdout)

		// Regex patterns for parsing bluetoothctl output
		// Handle ANSI escape codes that bluetoothctl may output
		deviceRegex := regexp.MustCompile(`\[(?:NEW|CHG)\] Device ([A-Fa-f0-9:]{17}) (.+)`)
		delDeviceRegex := regexp.MustCompile(`\[DEL\] Device ([A-Fa-f0-9:]{17})`)
		rssiRegex := regexp.MustCompile(`RSSI: (?:0x[a-fA-F0-9]+ )?\((-?\d+)\)`)
		nameRegex := regexp.MustCompile(`^([^R]+?)(?:\s+RSSI:|$)`)
		// Strip ANSI color codes and control characters
		ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[mK]|\r`)

		for scanner.Scan() {
			line := scanner.Text()
			// Clean line of ANSI escape codes and control characters
			cleanLine := ansiRegex.ReplaceAllString(line, "")
			cleanLine = strings.TrimSpace(cleanLine)

			// Handle device deletion
			if matches := delDeviceRegex.FindStringSubmatch(cleanLine); len(matches) >= 2 {
				macAddress := matches[1]
				ds.mutex.Lock()
				delete(ds.discoveredDevs, macAddress)
				ds.mutex.Unlock()
				continue
			}

			// Parse device discovery/change lines
			if matches := deviceRegex.FindStringSubmatch(cleanLine); len(matches) >= 3 {
				macAddress := matches[1]
				deviceInfo := matches[2]

				// Extract RSSI if present (handle both hex and decimal formats)
				rssi := 0
				if rssiMatches := rssiRegex.FindStringSubmatch(deviceInfo); len(rssiMatches) >= 2 {
					fmt.Sscanf(rssiMatches[1], "%d", &rssi)
				}

				// Extract device name using regex (everything before RSSI info)
				name := "Unknown Device"
				if nameMatches := nameRegex.FindStringSubmatch(deviceInfo); len(nameMatches) >= 2 {
					name = strings.TrimSpace(nameMatches[1])
				}
				if name == "" {
					name = "Unknown Device"
				}

				ds.mutex.Lock()
				// Update existing device or create new one
				if existing, exists := ds.discoveredDevs[macAddress]; exists {
					// Update RSSI and timestamp, keep other info
					existing.RSSI = rssi
					existing.Timestamp = time.Now()
					ds.discoveredDevs[macAddress] = existing
				} else {
					// Create new device
					ds.discoveredDevs[macAddress] = DiscoveredDevice{
						BluetoothDevice: BluetoothDevice{
							MacAddress: macAddress,
							Name:       name,
							RawLine:    cleanLine,
							Connected:  false,
							Paired:     false, // These are newly discovered
						},
						RSSI:      rssi,
						Timestamp: time.Now(),
					}
				}
				ds.mutex.Unlock()
			}
		}
	}()

	return nil
}

// StopDiscovery stops the scanning process
func (ds *DiscoveryScanner) StopDiscovery() error {
	if !ds.isScanning {
		return nil
	}

	if ds.cancel != nil {
		ds.cancel()
	}

	if ds.cmd != nil && ds.cmd.Process != nil {
		ds.cmd.Process.Kill()
		ds.cmd.Wait()
	}

	ds.isScanning = false
	return nil
}

// GetDiscoveredDevices returns the current list of discovered devices
func (ds *DiscoveryScanner) GetDiscoveredDevices() []DiscoveredDevice {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	devices := make([]DiscoveredDevice, 0, len(ds.discoveredDevs))
	for _, device := range ds.discoveredDevs {
		devices = append(devices, device)
	}
	return devices
}

// IsScanning returns whether discovery is currently active
func (ds *DiscoveryScanner) IsScanning() bool {
	return ds.isScanning
}

// ClearDiscoveredDevices clears the discovered devices list
func (ds *DiscoveryScanner) ClearDiscoveredDevices() {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	ds.discoveredDevs = make(map[string]DiscoveredDevice)
}

// DiscoveryTickCmd returns a command that sends periodic discovery updates
func DiscoveryTickCmd(scanner *DiscoveryScanner) tea.Cmd {
	return func() tea.Msg {
		if !scanner.IsScanning() {
			return nil
		}

		devices := scanner.GetDiscoveredDevices()
		return DiscoveryUpdateMsg{Devices: devices}
	}
}
