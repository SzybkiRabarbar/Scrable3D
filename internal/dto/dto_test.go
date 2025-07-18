package dto

import (
	"testing"
)

func TestValidatableImplementations(t *testing.T) {
	var _ Validatable = (*PlayData)(nil)
	var _ Validatable = (*ActionData)(nil)
	var _ Validatable = (*Char)(nil)
}

func TestValidParseID(t *testing.T) {
	testCases := []struct {
		name     string
		char     Char
		expected int64
	}{
		{
			name: "simple single digit",
			char: Char{
				Value:          "A",
				HtmlIdentifier: "char-A1",
			},
			expected: 1,
		},
		{
			name: "multiple digits",
			char: Char{
				Value:          "B",
				HtmlIdentifier: "char-B123",
			},
			expected: 123,
		},
		{
			name: "large number",
			char: Char{
				Value:          "Z",
				HtmlIdentifier: "char-Z999999",
			},
			expected: 999999,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := tc.char.ParseID()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if id != tc.expected {
				t.Errorf("expected id %d, got %d", tc.expected, id)
			}
		})
	}
}

func TestInvalidParseID(t *testing.T) {
	testCases := []struct {
		name        string
		char        Char
		expectedErr string
	}{
		{
			name: "missing char- prefix",
			char: Char{
				Value:          "A",
				HtmlIdentifier: "A1",
			},
			expectedErr: "char id must start with 'char-'",
		},
		{
			name: "mismatched letter",
			char: Char{
				Value:          "A",
				HtmlIdentifier: "char-B1",
			},
			expectedErr: "letter in char id (B) doesn't match its value (A)",
		},
		{
			name: "missing number",
			char: Char{
				Value:          "A",
				HtmlIdentifier: "char-A",
			},
			expectedErr: "invalid char id format",
		},
		{
			name: "invalid format with special characters",
			char: Char{
				Value:          "A",
				HtmlIdentifier: "char-A@1",
			},
			expectedErr: "invalid char id format",
		},
		{
			name: "lowercase letter",
			char: Char{
				Value:          "A",
				HtmlIdentifier: "char-a1",
			},
			expectedErr: "invalid char id format",
		},
		{
			name: "empty identifier",
			char: Char{
				Value:          "A",
				HtmlIdentifier: "",
			},
			expectedErr: "char id must start with 'char-'",
		},
		{
			name: "multiple letters",
			char: Char{
				Value:          "A",
				HtmlIdentifier: "char-AB1",
			},
			expectedErr: "invalid char id format",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.char.ParseID()
			if err == nil {
				t.Error("expected error, got nil")
				return
			}
			if err.Error() != tc.expectedErr {
				t.Errorf("expected error '%s', got '%s'", tc.expectedErr, err.Error())
			}
		})
	}
}
