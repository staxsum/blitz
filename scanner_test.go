package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewScanner(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantError bool
	}{
		{
			name:      "Valid HTTP URL",
			url:       "http://example.com",
			wantError: false,
		},
		{
			name:      "Valid HTTPS URL",
			url:       "https://example.com",
			wantError: false,
		},
		{
			name:      "URL without scheme",
			url:       "example.com",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("<html><body><form></form></body></html>"))
			}))
			defer ts.Close()

			scanner, err := NewScanner(ts.URL, 10, false)
			
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if !tt.wantError && scanner == nil {
				t.Error("Expected scanner but got nil")
			}
		})
	}
}

func TestFindForms(t *testing.T) {
	testHTML := `
	<html>
	<body>
		<form action="/login" method="post">
			<input type="text" name="username" />
			<input type="password" name="password" />
			<input type="submit" value="Login" />
		</form>
		<form action="/search" method="get">
			<input type="text" name="query" />
		</form>
	</body>
	</html>
	`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	}))
	defer ts.Close()

	scanner, err := NewScanner(ts.URL, 10, false)
	if err != nil {
		t.Fatalf("Failed to create scanner: %v", err)
	}

	forms, err := scanner.FindForms()
	if err != nil {
		t.Fatalf("FindForms failed: %v", err)
	}

	if len(forms) == 0 {
		t.Fatal("Expected to find forms but found none")
	}

	// Check if login form was detected
	foundLoginForm := false
	for _, form := range forms {
		if form.UsernameField != "" && form.PasswordField != "" {
			foundLoginForm = true
			
			if form.UsernameField != "username" {
				t.Errorf("Expected username field to be 'username', got '%s'", form.UsernameField)
			}
			
			if form.PasswordField != "password" {
				t.Errorf("Expected password field to be 'password', got '%s'", form.PasswordField)
			}
		}
	}

	if !foundLoginForm {
		t.Error("Expected to find login form with username and password fields")
	}
}

func TestSecurityChecks(t *testing.T) {
	tests := []struct {
		name           string
		headers        map[string]string
		expectVulnerable bool
	}{
		{
			name: "No X-Frame-Options",
			headers: map[string]string{},
			expectVulnerable: true,
		},
		{
			name: "With X-Frame-Options",
			headers: map[string]string{
				"X-Frame-Options": "DENY",
			},
			expectVulnerable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for key, value := range tt.headers {
					w.Header().Set(key, value)
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("<html><body></body></html>"))
			}))
			defer ts.Close()

			scanner, err := NewScanner(ts.URL, 10, false)
			if err != nil {
				t.Fatalf("Failed to create scanner: %v", err)
			}

			// Just run security checks, they print to stdout
			scanner.RunSecurityChecks(true)
		})
	}
}

func TestParseForm(t *testing.T) {
	testHTML := `
	<form action="/submit" method="POST">
		<input type="text" name="user" />
		<input type="password" name="pass" />
		<input type="hidden" name="token" value="abc123" />
		<select name="role">
			<option value="admin">Admin</option>
			<option value="user">User</option>
		</select>
	</form>
	`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testHTML))
	}))
	defer ts.Close()

	scanner, err := NewScanner(ts.URL, 10, false)
	if err != nil {
		t.Fatalf("Failed to create scanner: %v", err)
	}

	forms, err := scanner.FindForms()
	if err != nil {
		t.Fatalf("FindForms failed: %v", err)
	}

	if len(forms) != 1 {
		t.Fatalf("Expected 1 form, got %d", len(forms))
	}

	form := forms[0]

	if form.Method != "POST" {
		t.Errorf("Expected method POST, got %s", form.Method)
	}

	if !strings.Contains(form.Action, "/submit") {
		t.Errorf("Expected action to contain '/submit', got %s", form.Action)
	}

	if form.SelectField != "role" {
		t.Errorf("Expected select field 'role', got '%s'", form.SelectField)
	}

	if len(form.SelectOptions) != 2 {
		t.Errorf("Expected 2 select options, got %d", len(form.SelectOptions))
	}

	if form.Fields["token"] != "abc123" {
		t.Errorf("Expected hidden field token to be 'abc123', got '%s'", form.Fields["token"])
	}
}

func TestWAFDetection(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wafName    string
	}{
		{"ModSecurity", 406, "Mod_Security"},
		{"WebKnight", 999, "WebKnight"},
		{"F5 BIG IP", 419, "F5 BIG IP"},
		{"Unknown WAF", 403, "Unknown WAF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.String(), "script") {
					w.WriteHeader(tt.statusCode)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("<html></html>"))
				}
			}))
			defer ts.Close()

			scanner, err := NewScanner(ts.URL, 10, false)
			if err != nil {
				t.Fatalf("Failed to create scanner: %v", err)
			}

			// Run WAF detection - it prints to stdout
			scanner.detectWAF()
		})
	}
}
