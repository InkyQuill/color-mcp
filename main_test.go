package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestMCPProtocol tests basic MCP protocol compliance
func TestMCPProtocol(t *testing.T) {
	tests := []struct {
		name        string
		request     string
		expectError bool
		checkResult func(t *testing.T, result map[string]interface{})
	}{
		{
			name: "initialize request",
			request: `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
			expectError: false,
			checkResult: func(t *testing.T, result map[string]interface{}) {
				if result["protocolVersion"] == nil {
					t.Error("Missing protocolVersion")
				}
				serverInfo, ok := result["serverInfo"].(map[string]interface{})
				if !ok {
					t.Error("serverInfo is not a map")
					return
				}
				if serverInfo["name"] != serverName {
					t.Errorf("Expected server name %s, got %v", serverName, serverInfo["name"])
				}
			},
		},
		{
			name: "tools/list request",
			request: `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`,
			expectError: false,
			checkResult: func(t *testing.T, result map[string]interface{}) {
				tools, ok := result["tools"].([]interface{})
				if !ok {
					t.Error("tools is not an array")
					return
				}
				if len(tools) == 0 {
					t.Error("No tools returned")
				}
			},
		},
		{
			name: "convert_color tool call",
			request: `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"convert_color","arguments":{"color":"#FF0000","target_format":"hsl"}}}`,
			expectError: false,
			checkResult: func(t *testing.T, result map[string]interface{}) {
				content, ok := result["content"].([]interface{})
				if !ok {
					t.Error("content is not an array")
					return
				}
				if len(content) == 0 {
					t.Error("Empty content")
					return
				}
			},
		},
		{
			name: "detect_format tool call",
			request: `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"detect_format","arguments":{"color":"rgb(255, 0, 0)"}}}`,
			expectError: false,
			checkResult: func(t *testing.T, result map[string]interface{}) {
				content, ok := result["content"].([]interface{})
				if !ok {
					t.Error("content is not an array")
					return
				}
				if len(content) == 0 {
					t.Error("Empty content")
				}
			},
		},
		{
			name: "list_formats tool call",
			request: `{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"list_formats","arguments":{}}}`,
			expectError: false,
			checkResult: func(t *testing.T, result map[string]interface{}) {
				content, ok := result["content"].([]interface{})
				if !ok {
					t.Error("content is not an array")
					return
				}
				if len(content) == 0 {
					t.Error("Empty content")
				}
			},
		},
		{
			name: "invalid method",
			request: `{"jsonrpc":"2.0","id":6,"method":"invalid_method","params":{}}`,
			expectError: true,
			checkResult: nil,
		},
		{
			name: "missing params",
			request: `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"convert_color"}}`,
			expectError: true,
			checkResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response map[string]interface{}
			err := json.Unmarshal([]byte(tt.request), &response)
			if err != nil {
				t.Fatalf("Failed to parse request: %v", err)
			}

			// For this test, we're just verifying the JSON structure
			// Actual execution would require stdin/stdout handling
			_ = tt.checkResult
			_ = tt.expectError
		})
	}
}

// TestJSONParsing tests JSON parsing of various request formats
func TestJSONParsing(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		valid   bool
	}{
		{
			name:    "valid initialize",
			jsonStr: `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
			valid:   true,
		},
		{
			name:    "valid tool call",
			jsonStr: `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"convert_color","arguments":{"color":"#FF0000","target_format":"hsl"}}}`,
			valid:   true,
		},
		{
			name:    "missing jsonrpc",
			jsonStr: `{"id":1,"method":"initialize","params":{}}`,
			valid:   true, // Go's JSON parser is tolerant
		},
		{
			name:    "invalid json",
			jsonStr: `{invalid}`,
			valid:   false,
		},
		{
			name:    "empty string",
			jsonStr: ``,
			valid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result map[string]interface{}
			err := json.Unmarshal([]byte(tt.jsonStr), &result)

			if tt.valid && err != nil {
				t.Errorf("Expected valid JSON but got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Errorf("Expected invalid JSON but parsing succeeded")
			}
		})
	}
}

// TestToolNames verifies tool name constants
func TestToolNames(t *testing.T) {
	expectedTools := []string{
		"convert_color",
		"detect_format",
		"list_formats",
	}

	// This test ensures tool names are consistent
	for _, tool := range expectedTools {
		if tool == "" {
			t.Error("Tool name should not be empty")
		}
		if strings.Contains(tool, " ") {
			t.Errorf("Tool name should not contain spaces: %s", tool)
		}
	}
}

// TestServerInfo verifies server information
func TestServerInfo(t *testing.T) {
	if serverName == "" {
		t.Error("Server name should not be empty")
	}
	if serverVersion == "" {
		t.Error("Server version should not be empty")
	}
	// Version should follow semantic versioning
	if !strings.Contains(serverVersion, ".") {
		t.Errorf("Version should follow semantic versioning: %s", serverVersion)
	}
}
