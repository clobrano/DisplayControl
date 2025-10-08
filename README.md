# DisplayControl (`dispctl`)

A dynamic, TUI-based display control tool that allows you to manage monitor settings via DDC/CI commands. Built in Go with Bubbletea for an interactive terminal interface.

## Features

- 🎯 **Dynamic Configuration**: Auto-generate configs from your monitor's capabilities
- 📋 **Interactive TUI**: Navigate settings with an intuitive menu interface
- 🔍 **Dry-run Mode**: Preview commands before executing
- 🎨 **Breadcrumb Navigation**: Track your location in nested menus
- ⚡ **Single Binary**: No runtime dependencies except `ddcutil`
- ✅ **Error Handling**: Clear feedback on command success/failure

## Prerequisites

- Linux system with I2C support
- `ddcutil` installed and configured
- Go 1.23+ (for building from source)

### Installing ddcutil

```bash
# Debian/Ubuntu
sudo apt install ddcutil

# Fedora
sudo dnf install ddcutil

# Arch
sudo pacman -S ddcutil
```

### Setting up I2C permissions

Add your user to the `i2c` group:

```bash
sudo usermod -aG i2c $USER
```

Load the i2c-dev module:

```bash
sudo modprobe i2c-dev
echo "i2c-dev" | sudo tee /etc/modules-load.d/i2c-dev.conf
```

## Installation

### Building from Source

```bash
git clone https://github.com/clobrano/DisplayControl
cd DisplayControl
make build
sudo make install
```

## Usage

### Generate Configuration

Generate a configuration file from your monitor's capabilities:

```bash
dispctl --generate
```

This creates a configuration file at `~/.config/display-control/config.yaml`.

For multi-monitor setups, specify the display:

```bash
dispctl --generate --display 1
```

### Run the Interactive Menu

```bash
dispctl
```

Use a custom configuration file:

```bash
dispctl --file /path/to/config.yaml
```

### Dry-run Mode

Preview commands without executing them:

```bash
dispctl --dry-run
```

### Other Options

```bash
dispctl --version          # Show version
dispctl --help             # Show help
```

## Configuration File Format

The configuration file is generated from `ddcutil capabilities` output and follows this structure:

```yaml
brightness:
  _code: "0x10"
  _description: "Brightness"
  values:
    "0": "0x00"
    "50": "0x32"
    "100": "0x64"

input_source:
  _code: "0x60"
  _description: "Input Source"
  values:
    HDMI-1: "0x11"
    DisplayPort-1: "0x0f"
```

- `_code`: VCP feature code in hex format
- `_description`: Human-readable feature description
- `values`: Map of setting names to hex values

See `examples/sample-config.yaml` for a complete example.

## TUI Navigation

- **↑/↓** or **j/k**: Navigate menu items
- **Enter**: Select item / Execute command
- **b** or **Esc**: Go back one level
- **r**: Return to root menu
- **q** or **Ctrl+C**: Quit

## Command-Line Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--generate` | `-g` | Generate config from ddcutil capabilities |
| `--file PATH` | `-f` | Specify custom config file path |
| `--dry-run` | `-d` | Preview commands without executing |
| `--display ID` | | Specify display ID for ddcutil |
| `--version` | `-v` | Show version information |
| `--help` | `-h` | Show help message |

## Troubleshooting

### Permission Denied

If you get permission errors:

1. Ensure you're in the `i2c` group:
   ```bash
   groups $USER
   ```

2. Load the i2c-dev module:
   ```bash
   sudo modprobe i2c-dev
   ```

3. Log out and back in for group changes to take effect

### ddcutil Not Found

Install ddcutil for your distribution (see Prerequisites section).

### No Features in Configuration

Your monitor may not support DDC/CI or the feature list is empty. Check with:

```bash
ddcutil capabilities
```

### I2C Device Access

If you see "Cannot open I2C device":

1. Check that `/dev/i2c-*` devices exist
2. Verify i2c-dev module is loaded: `lsmod | grep i2c_dev`
3. Check permissions: `ls -l /dev/i2c-*`

## Development

### Running Tests

```bash
make test
```

### Building

```bash
make build
```

### Cleaning

```bash
make clean
```

## Project Structure

```
DisplayControl/
├── cmd/
│   └── dispctl/
│       └── main.go              # CLI entry point
├── internal/
│   ├── config/
│   │   ├── config.go            # YAML configuration parsing
│   │   └── config_test.go
│   ├── parser/
│   │   ├── capabilities.go      # ddcutil capabilities parser
│   │   └── capabilities_test.go
│   ├── generator/
│   │   ├── generator.go         # Configuration generator
│   │   └── generator_test.go
│   ├── ddcutil/
│   │   ├── executor.go          # DDC/CI command executor
│   │   └── executor_test.go
│   ├── ui/
│   │   ├── menu.go              # Bubbletea TUI
│   │   └── menu_test.go
│   └── version/
│       └── version.go           # Version information
├── examples/
│   └── sample-config.yaml       # Example configuration
├── Makefile
└── README.md
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) for the TUI
- Uses [ddcutil](https://www.ddcutil.com/) for DDC/CI communication
- Inspired by the need for better monitor control on Linux
