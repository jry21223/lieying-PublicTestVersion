package recon

import (
	"testing"
)

func TestExtractRootDomain(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"henu.edu.cn", "henu.edu.cn"},
		{"www.henu.edu.cn", "henu.edu.cn"},
		{"baidu.com", "baidu.com"},
		{"www.baidu.com", "baidu.com"},
		{"example.co.uk", "example.co.uk"},
		{"www.example.co.uk", "example.co.uk"},
		{"gov.cn", "gov.cn"},
		{"api.henu.edu.cn", "henu.edu.cn"},
	}

	for _, tt := range tests {
		result := ExtractRootDomain(tt.input)
		if result != tt.expected {
			t.Errorf("ExtractRootDomain(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}