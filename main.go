package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/InkyQuill/color-mcp/internal"
)

const (
	serverName    = "color-mcp"
	serverVersion = "1.0.0"
)

// MCP protocol structures
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type CallToolResult struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var req MCPRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			sendError(nil, -32700, "Parse error", err)
			continue
		}

		handleRequest(&req)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}

func handleRequest(req *MCPRequest) {
	// Notifications (requests without ID) should not receive responses
	isNotification := req.ID == nil

	switch req.Method {
	case "initialize":
		handleInitialize(req)
	case "tools/list":
		handleToolsList(req)
	case "tools/call":
		handleToolsCall(req)
	case "notifications/initialized":
		// Client notification that initialization is complete
		// No response needed for notifications
		break
	case "notifications/cancelled":
		// Request was cancelled
		// No response needed for notifications
		break
	case "notifications/progress":
		// Progress notification
		// No response needed for notifications
		break
	default:
		// Only send error response for requests, not notifications
		if !isNotification {
			sendError(req.ID, -32601, "Method not found", nil)
		}
	}
}

func handleInitialize(req *MCPRequest) {
	response := MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"serverInfo": map[string]string{
				"name":    serverName,
				"version": serverVersion,
			},
			"capabilities": map[string]interface{}{
				"tools": map[string]bool{},
			},
		},
	}
	sendResponse(response)
}

func handleToolsList(req *MCPRequest) {
	tools := []Tool{
		{
			Name:        "convert_color",
			Description: "Convert colors between different web color formats (HEX, RGB, HSL, OKLCH, LAB, XYZ, HWB, CMYK, etc.)",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"color": {
						Type:        "string",
						Description: "Input color value in any supported format (e.g., '#FF0000', 'rgb(255, 0, 0)', 'hsl(0, 100%, 50%)')",
					},
					"target_format": {
						Type:        "string",
						Description: "Target color format",
						Enum:        internal.GetSupportedFormats(),
					},
					"preserve_alpha": {
						Type:        "boolean",
						Description: "Whether to preserve the alpha channel (default: true)",
					},
				},
				Required: []string{"color", "target_format"},
			},
		},
		{
			Name:        "detect_format",
			Description: "Detect the format of an input color string",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"color": {
						Type:        "string",
						Description: "Color value to detect format from",
					},
				},
				Required: []string{"color"},
			},
		},
		{
			Name:        "list_formats",
			Description: "List all supported color formats",
			InputSchema: InputSchema{
				Type:       "object",
				Properties: map[string]Property{},
				Required:   []string{},
			},
		},
		{
			Name:        "compare_colors",
			Description: "Compare two colors for perceptual similarity, contrast ratio, and component differences",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"color1": {
						Type:        "string",
						Description: "First color value in any supported format (e.g., '#FF0000', 'rgb(255, 0, 0)', 'hsl(0, 100%, 50%)')",
					},
					"color2": {
						Type:        "string",
						Description: "Second color value in any supported format",
					},
					"detailed": {
						Type:        "boolean",
						Description: "Whether to include detailed component breakdown (default: false)",
					},
				},
				Required: []string{"color1", "color2"},
			},
		},
	}

	response := MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"tools": tools,
		},
	}
	sendResponse(response)
}

func handleToolsCall(req *MCPRequest) {
	var params ToolCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		sendError(req.ID, -32602, "Invalid params", err)
		return
	}

	var result CallToolResult
	var err error

	switch params.Name {
	case "convert_color":
		result, err = convertColor(params.Arguments)
	case "detect_format":
		result, err = detectFormat(params.Arguments)
	case "list_formats":
		result, err = listFormats(params.Arguments)
	case "compare_colors":
		result, err = compareColors(params.Arguments)
	default:
		sendError(req.ID, -32601, "Unknown tool: "+params.Name, nil)
		return
	}

	if err != nil {
		result = CallToolResult{
			Content: []ContentItem{
				{Type: "text", Text: fmt.Sprintf("Error: %v", err)},
			},
			IsError: true,
		}
	}

	response := MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
	sendResponse(response)
}

func convertColor(args map[string]interface{}) (CallToolResult, error) {
	color, ok := args["color"].(string)
	if !ok {
		return CallToolResult{}, fmt.Errorf("color parameter is required and must be a string")
	}

	targetFormat, ok := args["target_format"].(string)
	if !ok {
		return CallToolResult{}, fmt.Errorf("target_format parameter is required and must be a string")
	}

	preserveAlpha := true
	if pa, ok := args["preserve_alpha"].(bool); ok {
		preserveAlpha = pa
	}

	// Detect input format first
	inputFormat, err := internal.DetectInputFormat(color)
	if err != nil {
		return CallToolResult{}, fmt.Errorf("failed to detect input format: %w", err)
	}

	// Convert
	output, err := internal.Convert(color, targetFormat, preserveAlpha)
	if err != nil {
		return CallToolResult{}, err
	}

	// Format result
	resultText := fmt.Sprintf("Input color: %s (format: %s)\nOutput color: %s (format: %s)\nAlpha preserved: %t",
		color, inputFormat, output, targetFormat, preserveAlpha)

	return CallToolResult{
		Content: []ContentItem{
			{Type: "text", Text: resultText},
		},
	}, nil
}

func detectFormat(args map[string]interface{}) (CallToolResult, error) {
	color, ok := args["color"].(string)
	if !ok {
		return CallToolResult{}, fmt.Errorf("color parameter is required and must be a string")
	}

	format, err := internal.DetectInputFormat(color)
	if err != nil {
		return CallToolResult{}, err
	}

	resultText := fmt.Sprintf("Color: %s\nDetected format: %s", color, format)

	return CallToolResult{
		Content: []ContentItem{
			{Type: "text", Text: resultText},
		},
	}, nil
}

func listFormats(args map[string]interface{}) (CallToolResult, error) {
	formats := internal.GetSupportedFormats()
	resultText := "Supported color formats:\n" + strings.Join(formats, ", ")

	return CallToolResult{
		Content: []ContentItem{
			{Type: "text", Text: resultText},
		},
	}, nil
}

func compareColors(args map[string]interface{}) (CallToolResult, error) {
	color1, ok := args["color1"].(string)
	if !ok {
		return CallToolResult{}, fmt.Errorf("color1 parameter is required and must be a string")
	}

	color2, ok := args["color2"].(string)
	if !ok {
		return CallToolResult{}, fmt.Errorf("color2 parameter is required and must be a string")
	}

	detailed := false
	if d, ok := args["detailed"].(bool); ok {
		detailed = d
	}

	result, err := internal.CompareColors(color1, color2)
	if err != nil {
		return CallToolResult{}, err
	}

	var resultText string
	if detailed {
		resultText = internal.FormatComparisonDetailed(result)
	} else {
		resultText = internal.FormatComparisonBasic(result)
	}

	return CallToolResult{
		Content: []ContentItem{
			{Type: "text", Text: resultText},
		},
	}, nil
}

func sendResponse(resp MCPResponse) {
	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

func sendError(id interface{}, code int, message string, err error) {
	resp := MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}
	if err != nil {
		resp.Error.Message += fmt.Sprintf(": %v", err)
	}
	sendResponse(resp)
}
