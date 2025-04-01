package configtype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "valid base64 string",
			input:    "SGVsbG8gV29ybGQ=",
			expected: "Hello World",
			wantErr:  false,
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "invalid base64 string",
			input:    "not-base64!",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "base64 with padding",
			input:    "SGVsbG8=",
			expected: "Hello",
			wantErr:  false,
		},
		{
			name:     "base64 with multiple padding",
			input:    "SGVsbG8gV29ybGQ==",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Base64
			err := b.UnmarshalText([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(b))
		})
	}
}
