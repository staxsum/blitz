package main

import (
	"testing"
)

func TestFormDescription(t *testing.T) {
	form := &Form{
		Action:        "/login",
		UsernameField: "user",
		PasswordField: "pass",
		SelectField:   "role",
	}

	desc := form.Description()

	if desc == "" {
		t.Error("Expected non-empty description")
	}

	expectedParts := []string{"/login", "user", "pass", "role"}
	for _, part := range expectedParts {
		if !contains(desc, part) {
			t.Errorf("Expected description to contain '%s', got: %s", part, desc)
		}
	}
}

func TestIsLoginForm(t *testing.T) {
	tests := []struct {
		name     string
		form     *Form
		expected bool
	}{
		{
			name: "Valid login form",
			form: &Form{
				UsernameField: "username",
				PasswordField: "password",
			},
			expected: true,
		},
		{
			name: "Only username field",
			form: &Form{
				UsernameField: "username",
			},
			expected: false,
		},
		{
			name: "Only password field",
			form: &Form{
				PasswordField: "password",
			},
			expected: false,
		},
		{
			name:     "No fields",
			form:     &Form{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.form.IsLoginForm()
			if result != tt.expected {
				t.Errorf("IsLoginForm() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHasPasswordField(t *testing.T) {
	tests := []struct {
		name     string
		form     *Form
		expected bool
	}{
		{
			name: "With password field",
			form: &Form{
				PasswordField: "password",
			},
			expected: true,
		},
		{
			name:     "Without password field",
			form:     &Form{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.form.HasPasswordField()
			if result != tt.expected {
				t.Errorf("HasPasswordField() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[0:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
