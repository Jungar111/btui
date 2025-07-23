# btui

A Terminal User Interface (TUI) for managing Bluetooth devices through bluetoothctl.

## Features

- **Interactive Device Listing**: Browse available Bluetooth devices with connection status
- **Device Connection**: Select and connect to Bluetooth devices with real-time feedback
- **Clean Terminal UI**: Built with the Charm stack for a polished terminal experience

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

# Build the application
go build -o btui .

# Run btui
./btui --help
```

## Usage

### List Devices
View all available Bluetooth devices:
```bash
./btui list-devices
```

### Connect to Device
Select and connect to a Bluetooth device:
```bash
./btui connect
```

## Commands

- `list-devices` - List and select Bluetooth devices
- `connect` - Connect to a Bluetooth device

## Requirements

- Go 1.24 or later
- `bluetoothctl` command available in PATH (install with `sudo apt install bluez-utils` on Ubuntu/Debian or `sudo dnf install bluez` on Fedora)
- Linux system with Bluetooth support

## Architecture

btui follows a modular architecture with separate commands for different Bluetooth operations:

- `cmd/` - Individual command implementations
- `internal/bluetooth/` - Shared Bluetooth utilities and device management
- `internal/ui/` - Common UI components and styling

## Development

```bash
# Build
go build -o btui .

# Run tests
go test ./...

# Format code
go fmt ./...

# Lint code
go vet ./...
```

## Technologies

- **Bubble Tea** - TUI framework
- **Bubbles** - Pre-built TUI components  
- **Lipgloss** - Terminal styling
- **Cobra** - CLI framework

## License

See LICENSE file for details.