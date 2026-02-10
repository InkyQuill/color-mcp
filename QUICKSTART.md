# Quick Start Guide

## Installation

```bash
# Build
go build -o color-mcp

# Or install directly
go install github.com/inky/color-mcp@latest
```

## MCP Configuration

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "color-converter": {
      "command": "/path/to/color-mcp"
    }
  }
}
```

## Usage Examples

### Convert HEX to HSL
```
Convert #FF5733 to HSL
```

### Convert RGB to OKLCH
```
Convert rgb(255, 0, 0) to oklch
```

### Detect color format
```
Detect the format of hsl(120, 100%, 50%)
```

### List all formats
```
List all supported color formats
```

## Supported Formats

- `hex` - #FF0000, #FF000080
- `rgb` / `rgba` - rgb(255, 0, 0), rgba(255, 0, 0, 0.5)
- `hsl` / `hsla` - hsl(0, 100%, 50%), hsla(0, 100%, 50%, 0.5)
- `hsb` / `hsv` - hsb(0, 100%, 100%)
- `oklch` - oklch(0.5 0.1 120)
- `lab` - lab(50 50 50)
- `xyz` - xyz(0.5 0.5 0.5)
- `hwb` - hwb(0 0% 0%)
- `cmyk` - cmyk(0% 100% 100% 0%)

## Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```
