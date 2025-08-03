# btui

A Terminal User Interface (TUI) for managing Bluetooth devices through bluetoothctl.

## Features

- **🔍 Real-time Device Discovery**: Find nearby Bluetooth devices that aren't paired yet
- **📡 Live Signal Strength**: See RSSI values and device signal strength in real-time
- **📱 Smart Device Organization**: Three-tier sorting (Connected → Paired → Discovered)
- **🎨 Terminal-Adaptive Interface**: Colors respect your terminal theme using Charm stack
- **⚡ Direct Launch**: Opens directly to scanning interface for immediate productivity
- **⌨️ Intuitive Controls**: Keyboard shortcuts for all major actions
- **🛡️ Robust Parsing**: Handles bluetoothctl's ANSI colors and real-time output
- **🔄 Position Preservation**: List maintains position during real-time updates

## Installation

### Quick Install (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd btui

# Run the install script (builds and installs to /usr/local/bin)
./install.sh
```

### Manual Installation

```bash
# Clone the repository
git clone <repository-url>
cd btui

# Build the application (using justfile)
just build

# Or build manually
go build -o dist/btui .

# Run btui
./dist/btui --help
```

## Usage

### Interactive Interface (Recommended)
Launch the interactive scan interface directly:
```bash
btui
```
The application launches directly into the scanning interface, providing immediate access to all Bluetooth device management features.

### Direct Command Access

#### Scan for Devices
Discover nearby Bluetooth devices in real-time, including devices that aren't paired yet:
```bash
btui scan
```

**Real-time Discovery Features:**
- **Live device discovery** - Finds nearby devices as they become available
- **RSSI signal strength** - Shows signal strength (e.g., -72 dBm) 
- **Device lifecycle** - Automatically updates as devices appear/disappear
- **Mixed device view** - Shows both paired and newly discovered devices

**Scan Controls:**
- `s` - Start/stop real-time discovery
- `c` - Connect to selected device (works with both paired and discovered)
- `d` - Disconnect from selected device
- `r` - Refresh paired device list
- `↑/↓` - Navigate device list
- `q` - Quit

#### List Paired Devices
View all paired Bluetooth devices:
```bash
btui list-devices
```

#### Connect to Device
Select and connect to a paired Bluetooth device:
```bash
btui connect
```

#### Disconnect from Device
Select and disconnect from a connected Bluetooth device:
```bash
btui disconnect
```

## Device Status Display

Devices are organized in a prioritized list with colored status indicators:

**Device Priority:**
1. **Connected** - Currently active connections (green)
2. **Paired** - Previously paired devices (yellow)
3. **Discovered** - Newly found devices during scanning (cyan)

**Status Information:**
- **Signal Strength** - RSSI values shown for discovered devices (e.g., "RSSI: -72")
- **MAC Addresses** - Device hardware addresses displayed in muted text
- **Real-time Updates** - Live updates as devices appear, change, or disappear
- **Smart Sorting** - Alphabetical sorting within each status category

## Commands

- `scan` - **Real-time discovery** of nearby Bluetooth devices (both paired and unpaired)
- `list-devices` - List and select paired Bluetooth devices only
- `connect` - Connect to a paired Bluetooth device
- `disconnect` - Disconnect from a connected Bluetooth device

## Requirements

- Go 1.24 or later
- `bluetoothctl` command available in PATH (install with `sudo apt install bluez-utils` on Ubuntu/Debian or `sudo dnf install bluez` on Fedora)
- Linux system with Bluetooth support

## Architecture

btui follows a modular architecture with separate commands for different Bluetooth operations:

- **`cmd/`** - Individual command implementations
  - `listdevices/` - Device listing functionality
  - `scan/` - Real-time scanning and discovery
  - `connect/` - Device connection interface
  - `disconnect/` - Device disconnection interface
  - `root.go` - Root command and CLI setup
- **`internal/bluetooth/`** - Shared Bluetooth utilities and device management
  - `commands.go` - Bluetooth command implementations (connect, disconnect)
  - `scanner.go` - Paired device scanning and parsing logic
  - `discovery.go` - **Real-time device discovery engine** (NEW)
  - `types.go` - Bluetooth device data structures
  - `scanner_test.go` - Comprehensive test suite
- **`internal/ui/`** - Common UI components and styling
  - `list.go` - Generic list component
  - `styling.go` - Centralized styling definitions
- **`internal/menu/`** - Main menu interface
  - Navigation and sub-program management

## Development

### Using Justfile (Recommended)

```bash
# Build the application
just build

# Run tests
just test

# Format code
just fmt

# Vet code for issues
just vet

# Run all checks (fmt, vet, test)
just check

# Full pipeline (check + build)
just all

# Clean build artifacts
just clean
```

### Manual Commands

```bash
# Build
go build -o dist/btui .

# Run tests
go test ./...

# Test real-time discovery (standalone)
go run discovery_test_standalone.go

# Format code
go fmt ./...

# Lint code
go vet ./...

# Tidy dependencies
go mod tidy
```

### Testing Real Device Discovery

To test the real-time discovery functionality:

```bash
# Run the standalone discovery test
go run discovery_test_standalone.go

# Or test via the main application
./dist/btui scan
# Press 's' to start discovery and watch for new devices
```

## Technologies

- **Bubble Tea** - TUI framework
- **Bubbles** - Pre-built TUI components  
- **Lipgloss** - Terminal styling
- **Cobra** - CLI framework

## License

See LICENSE file for details.