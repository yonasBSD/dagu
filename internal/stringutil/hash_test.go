package stringutil

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase58EncodeSHA256(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string // We'll validate format rather than exact value
	}{
		{
			name:  "simple string",
			input: "hello",
		},
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "child DAG ID format",
			input: "12345:process-data:{\"REGION\":\"us-east-1\",\"VERSION\":\"1.0.0\"}",
		},
		{
			name:  "unicode characters",
			input: "Hello 世界 🌍",
		},
		{
			name:  "long input",
			input: "parent-dag-run-id-12345:parallel-step-name:{\"param1\":\"value1\",\"param2\":\"value2\",\"param3\":\"value3\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base58EncodeSHA256(tt.input)

			// Validate result is non-empty
			assert.NotEmpty(t, result, "base58 encoded hash should not be empty")

			// Validate result contains only base58 characters
			for _, c := range result {
				assert.Contains(t, base58Alphabet, string(c), "result should only contain base58 characters")
			}

			// Validate consistency - same input should always produce same output
			result2 := Base58EncodeSHA256(tt.input)
			assert.Equal(t, result, result2, "same input should produce same output")
		})
	}
}

func TestBase58Encode(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "empty bytes",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "single zero byte",
			input:    []byte{0},
			expected: "1",
		},
		{
			name:     "multiple zero bytes",
			input:    []byte{0, 0, 0},
			expected: "111",
		},
		{
			name:     "simple bytes",
			input:    []byte{1, 2, 3},
			expected: "Ldp",
		},
		{
			name:     "SHA256 hash",
			input:    func() []byte { h := sha256.Sum256([]byte("test")); return h[:] }(),
			expected: "Bjj4AWTNrjQVHqgWbP2XaxXz4DYH1WZMyERHxsad7b2w",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base58Encode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBase58EncodeSHA256_ChildDAGScenarios(t *testing.T) {
	// Test specific scenarios for child DAG ID generation
	tests := []struct {
		name           string
		parentRunID    string
		stepName       string
		params         string
		expectedLength int // Expected length range
	}{
		{
			name:           "simple child DAG",
			parentRunID:    "parent-12345",
			stepName:       "process",
			params:         `{"env":"prod"}`,
			expectedLength: 40, // Base58 encoded SHA256 is typically 43-44 chars
		},
		{
			name:           "child DAG with complex params",
			parentRunID:    "workflow-abc-123",
			stepName:       "etl-pipeline",
			params:         `{"AWS_REGION":"us-east-1","BATCH_SIZE":"1000","MODE":"parallel"}`,
			expectedLength: 40,
		},
		{
			name:           "nested child DAG scenario",
			parentRunID:    "root-workflow:child-workflow:grandchild-12345",
			stepName:       "data-processor",
			params:         `{"input":"/data/raw","output":"/data/processed"}`,
			expectedLength: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Construct the input as it would be in real usage
			input := fmt.Sprintf("%s:%s:%s", tt.parentRunID, tt.stepName, tt.params)

			result := Base58EncodeSHA256(input)

			// Validate result properties
			assert.NotEmpty(t, result)
			assert.GreaterOrEqual(t, len(result), tt.expectedLength)
			assert.LessOrEqual(t, len(result), tt.expectedLength+4) // Allow some variance

			// Ensure different inputs produce different outputs
			altInput := input + "-modified"
			altResult := Base58EncodeSHA256(altInput)
			assert.NotEqual(t, result, altResult, "different inputs should produce different hashes")

			// Validate no ambiguous characters
			for _, c := range result {
				assert.NotContains(t, "0OlI", string(c), "result should not contain ambiguous characters")
			}
		})
	}
}

func BenchmarkBase58EncodeSHA256(b *testing.B) {
	input := "parent-dag-run-12345:step-name:{\"key\":\"value\"}"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Base58EncodeSHA256(input)
	}
}

func BenchmarkBase58Encode(b *testing.B) {
	h := sha256.Sum256([]byte("benchmark test data"))
	data := h[:]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Base58Encode(data)
	}
}
