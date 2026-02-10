# Color MCP Server

A Model Context Protocol (MCP) server for converting between various web color formats. This server provides tools to detect color formats and convert colors between HEX, RGB, HSL, HSB/HSV, OKLCH, LAB, XYZ, HWB, and CMYK formats.

## Features

- **Auto-detection**: Automatically detects input color format
- **Multiple formats**: Supports 10+ color formats
- **Alpha channel**: Preserves or strips alpha channel as needed
- **Color comparison**: Perceptual similarity analysis using OKLCH ΔE
- **Accessibility**: WCAG contrast ratio calculations
- **High precision**: Uses accurate color space conversions
- **Fast**: Pure Go implementation with no external dependencies
- **Comprehensive tests**: >90% test coverage

## Supported Color Formats

| Format | Examples | Description |
|--------|----------|-------------|
| HEX | `#FF0000`, `#FF000080` | Hexadecimal RGB/RGBA |
| RGB | `rgb(255, 0, 0)`, `rgb(100%, 0%, 0%)` | RGB color space |
| RGBA | `rgba(255, 0, 0, 0.5)` | RGB with alpha |
| HSL | `hsl(0, 100%, 50%)` | Hue, Saturation, Lightness |
| HSLA | `hsla(0, 100%, 50%, 0.5)` | HSL with alpha |
| HSB/HSV | `hsb(0, 100%, 100%)` | Hue, Saturation, Brightness/Value |
| OKLCH | `oklch(0.5 0.1 120)` | Perceptually uniform color space |
| LAB | `lab(50 50 50)` | CIE LAB color space |
| XYZ | `xyz(0.5 0.5 0.5)` | CIE XYZ color space |
| HWB | `hwb(0 0% 0%)` | Hue, Whiteness, Blackness |
| CMYK | `cmyk(0% 100% 100% 0%)` | Cyan, Magenta, Yellow, Key (Black) |

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/inky/color-mcp.git
cd color-mcp

# Build the binary
go build -o color-mcp

# (Optional) Install to system
sudo install color-mcp /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/inky/color-mcp@latest
```

## MCP Configuration

Add to your MCP client configuration (typically in `~/.config/claude/mcp.json` or similar):

```json
{
  "mcpServers": {
    "color-converter": {
      "command": "/path/to/color-mcp"
    }
  }
}
```

**Replace** `/path/to/color-mcp` with the actual path to the binary (e.g., `/usr/local/bin/color-mcp` or `~/go/bin/color-mcp`).

## Usage

### Available Tools

#### 1. convert_color

Convert a color from one format to another.

**Parameters:**
- `color` (string, required): Input color value in any supported format
- `target_format` (string, required): Target format (hex, rgb, hsl, hsla, hsb, oklch, lab, xyz, hwb, cmyk)
- `preserve_alpha` (boolean, optional): Whether to preserve alpha channel (default: true)

**Example:**
```
Convert #FF0000 to HSL format
```

#### 2. detect_format

Detect the format of an input color string.

**Parameters:**
- `color` (string, required): Color value to detect format from

**Example:**
```
Detect the format of rgb(255, 0, 0)
```

#### 3. list_formats

List all supported color formats.

**Example:**
```
List all supported color formats
```

#### 4. compare_colors

Compare two colors for perceptual similarity, contrast ratio, and component differences.

**Parameters:**
- `color1` (string, required): First color value in any supported format
- `color2` (string, required): Second color value in any supported format
- `detailed` (boolean, optional): Whether to include detailed component breakdown (default: false)

**Example:**
```
Compare #FF0000 and #00FF00 with detailed output
```

Result (basic):
```
Color Comparison: #FF0000 vs #00FF00
Perceptual Difference: 0.520 ΔE
Verdict: different
Contrast Ratio: 2.91:1 (Fail)
```

Result (detailed):
```
Color Comparison: #FF0000 (hex) vs #00FF00 (hex)

Perceptual Difference: 0.520 ΔE
Verdict: different

Component Breakdown:
  Hue Difference: 120.0°
  Lightness Difference: 0.0%
  Saturation Difference: 0.0%

Contrast Ratio: 2.91:1
WCAG Grade: Fail
```

## Examples

### Converting HEX to HSL

```
Convert #FF5733 to hsl
```

Result:
```
Input color: #FF5733 (format: hex)
Output color: hsl(11, 100%, 60%) (format: hsl)
Alpha preserved: true
```

### Converting RGB to OKLCH

```
Convert rgb(255, 0, 0) to oklch
```

Result:
```
Input color: rgb(255, 0, 0) (format: rgb)
Output color: oklch(0.6280 0.2577 29.23 / 1.00) (format: oklch)
Alpha preserved: true
```

### Detecting Color Format

```
Detect the format of hsl(120, 100%, 50%)
```

Result:
```
Color: hsl(120, 100%, 50%)
Detected format: hsl
```

### Alpha Channel Handling

```
Convert rgba(255, 0, 0, 0.5) to hex
```

Result:
```
Input color: rgba(255, 0, 0, 0.5) (format: rgba)
Output color: #FF000080 (format: hex)
Alpha preserved: true
```

```
Convert rgba(255, 0, 0, 0.5) to rgb without preserving alpha
```

Result:
```
Input color: rgba(255, 0, 0, 0.5) (format: rgba)
Output color: rgb(255, 0, 0) (format: rgb)
Alpha preserved: false
```

### Color Comparison

```
Compare #000000 and #FFFFFF
```

Result:
```
Color Comparison: #000000 vs #FFFFFF
Perceptual Difference: 1.000 ΔE
Verdict: different
Contrast Ratio: 21.00:1 (AAA)
```

```
Compare hsl(350, 100%, 50%) and hsl(10, 100%, 50%) with detailed output
```

Result:
```
Color Comparison: hsl(350, 100%, 50%) (hsl) vs hsl(10, 100%, 50%) (hsl)

Perceptual Difference: 0.032 ΔE
Verdict: slightly different

Component Breakdown:
  Hue Difference: 20.0°
  Lightness Difference: 0.0%
  Saturation Difference: 0.0%

Contrast Ratio: 1.06:1
WCAG Grade: Fail
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...
```

### Building

```bash
# Build for current platform
go build

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o color-mcp-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o color-mcp-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o color-mcp-darwin-arm64
GOOS=windows GOARCH=amd64 go build -o color-mcp-windows-amd64.exe
```

### Project Structure

```
color-mcp/
├── internal/
│   ├── types.go       # Color format detection and parsing
│   ├── convert.go     # Color conversion algorithms
│   ├── converter.go   # Main conversion logic and formatting
│   ├── compare.go     # Color comparison and contrast calculation
│   ├── constants.go   # Color space constants and thresholds
│   ├── value_objects.go   # Channel value types
│   └── *_test.go      # Comprehensive tests
├── main.go            # MCP server implementation
├── go.mod
├── go.sum
└── README.md
```

## Color Conversion Algorithm

The converter uses OKLCH as the intermediate format for highest quality conversions:

1. **Input Detection**: Parse and detect the input format
2. **RGB Conversion**: Convert input format to RGB
3. **Target Conversion**: Convert RGB to target format
4. **Formatting**: Format output according to target specification

Special handling:
- **OKLCH**: Uses proper CIE XYZ intermediate for perceptual accuracy
- **LAB/XYZ**: Standard CIE color space conversions
- **CMYK**: Applies proper black key generation
- **Alpha**: Preserved or stripped based on parameter

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Author

Created by Inky (@inky)

## Acknowledgments

- Color space conversion formulas based on CSS Color Module Level 4
- OKLCH implementation based on the work by Björn Ottosson
- MCP protocol specification by Anthropic
