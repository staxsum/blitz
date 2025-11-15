package main

import (
	"fmt"
	"strings"
)

type Form struct {
	Index         int
	Action        string
	Method        string
	UsernameField string
	PasswordField string
	SelectField   string
	SelectOptions []string
	Fields        map[string]string
}

func (f *Form) Description() string {
	parts := []string{fmt.Sprintf("Action: %s", f.Action)}
	
	if f.UsernameField != "" {
		parts = append(parts, fmt.Sprintf("Username: %s", f.UsernameField))
	}
	
	if f.PasswordField != "" {
		parts = append(parts, fmt.Sprintf("Password: %s", f.PasswordField))
	}
	
	if f.SelectField != "" {
		parts = append(parts, fmt.Sprintf("Select: %s", f.SelectField))
	}
	
	return strings.Join(parts, ", ")
}

func (f *Form) IsLoginForm() bool {
	return f.UsernameField != "" && f.PasswordField != ""
}

func (f *Form) HasPasswordField() bool {
	return f.PasswordField != ""
}
