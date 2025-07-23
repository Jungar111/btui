# btui

A Terminal User Interface (TUI) for managing Bluetooth devices through bluetoothctl.

## Features

- **🔍 Real-time Device Discovery**: Find nearby Bluetooth devices that aren't paired yet
- **📡 Live Signal Strength**: See RSSI values and device signal strength in real-time
- **📱 Interactive Device Management**: Browse, connect, and disconnect from Bluetooth devices
- **🔗 Connection Status Tracking**: Visual indicators for discovered (🔍), paired (📱), and connected (🔗) devices
- **⚡ Multiple Interface Options**: Main menu interface or direct command access
- **🎨 Clean Terminal UI**: Built with the Charm stack for a polished terminal experience
- **⌨️ Intuitive Controls**: Keyboard shortcuts for all major actions
- **🛡️ Robust Parsing**: Handles bluetoothctl's ANSI colors and real-time output

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

### Main Menu Interface (Recommended)
Launch the interactive main menu:
```bash
btui
```
Navigate through options:
- **List Devices** - Browse paired Bluetooth devices
- **Scan & Connect** - Discover and connect to nearby devices
- **Connect** - Connect to a paired device
- **Disconnect** - Disconnect from a connected device

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

## Device Status Indicators

- 🔍 **Discovered** - Newly found device during real-time scanning (not paired yet)
- 📱 **Paired** - Previously paired device (appears in `bluetoothctl devices`)
- 🔗 **Connected** - Currently connected device

**Additional Information:**
- **RSSI values** shown for discovered devices (e.g., "RSSI: -72")
- **Device names** extracted from Bluetooth advertisements when available
- **Real-time updates** as devices move in/out of range

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