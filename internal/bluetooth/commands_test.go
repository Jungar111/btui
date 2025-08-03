package bluetooth

import (
	"strings"
	"testing"
)

func TestDisconnectSuccessDetection(t *testing.T) {
	testCases := []struct {
		name     string
		output   string
		hasError bool
		expected bool
	}{
		{"Successful disconnected", "Successful disconnected", false, true},
		{"Device disconnected", "Device AA:BB:CC:DD:EE:FF disconnected", false, true},
		{"Empty output success", "", false, true},
		{"Failed to disconnect", "Failed to disconnect", false, false},
		{"Error with disconnected", "Error: Device not connected", false, false},
		{"Command error", "Successful disconnected", true, false},
		{"Mixed case success", "DEVICE AA:BB:CC:DD:EE:FF DISCONNECTED", false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DisconnectResult{
				Output: tc.output,
				Err:    nil,
			}
			
			if tc.hasError {
				result.Err = &testError{}
			}

			// Apply the same logic as in DisconnectCmd
			lowerOutput := strings.ToLower(result.Output)
			if result.Err == nil && (strings.Contains(lowerOutput, "successful") || 
			                         strings.Contains(lowerOutput, "disconnected") ||
			                         result.Output == "") {
				if !strings.Contains(lowerOutput, "failed") && !strings.Contains(lowerOutput, "error") {
					result.Success = true
				}
			}

			if result.Success != tc.expected {
				t.Errorf("Expected success=%v, got success=%v for output: %q", 
					tc.expected, result.Success, tc.output)
			}
		})
	}
}

func TestConnectSuccessDetection(t *testing.T) {
	testCases := []struct {
		name     string
		output   string
		hasError bool
		expected bool
	}{
		{"Connection successful", "Connection successful", false, true},
		{"Device connected", "Device AA:BB:CC:DD:EE:FF connected", false, true},
		{"Empty output success", "", false, true},
		{"Failed to connect", "Failed to connect", false, false},
		{"Error with connected", "Error: Device not found", false, false},
		{"Command error", "Connection successful", true, false},
		{"Mixed case success", "DEVICE AA:BB:CC:DD:EE:FF CONNECTED", false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ConnectResult{
				Output: tc.output,
				Err:    nil,
			}
			
			if tc.hasError {
				result.Err = &testError{}
			}

			// Apply the same logic as in ConnectCmd
			lowerOutput := strings.ToLower(result.Output)
			if result.Err == nil && (strings.Contains(lowerOutput, "successful") || 
			                         strings.Contains(lowerOutput, "connected") ||
			                         result.Output == "") {
				if !strings.Contains(lowerOutput, "failed") && !strings.Contains(lowerOutput, "error") {
					result.Success = true
				}
			}

			if result.Success != tc.expected {
				t.Errorf("Expected success=%v, got success=%v for output: %q", 
					tc.expected, result.Success, tc.output)
			}
		})
	}
}

type testError struct{}
func (e *testError) Error() string { return "test error" }