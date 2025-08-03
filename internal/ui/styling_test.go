package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestStylesAreDefined(t *testing.T) {
	// Test that all styles are properly defined and not nil
	styles := []struct {
		name  string
		style lipgloss.Style
	}{
		{"TitleStyle", TitleStyle},
		{"SubtitleStyle", SubtitleStyle},
		{"HelpStyle", HelpStyle},
		{"ItemStyle", ItemStyle},
		{"SelectedItemStyle", SelectedItemStyle},
		{"ConnectedItemStyle", ConnectedItemStyle},
		{"ConnectedSelectedItemStyle", ConnectedSelectedItemStyle},
		{"PaginationStyle", PaginationStyle},
	}

	for _, test := range styles {
		// Test that style can render text without panicking
		rendered := test.style.Render("test")
		if rendered == "" {
			t.Errorf("%s produced empty render", test.name)
		}
	}
}

func TestSuccessStyle(t *testing.T) {
	style := SuccessStyle()
	rendered := style.Render("Success message")

	if rendered == "" {
		t.Error("SuccessStyle produced empty render")
	}

	// Success style should contain the text even if formatting is not visible in tests
	if !strings.Contains(rendered, "Success message") {
		t.Error("SuccessStyle should contain the original text")
	}
}

func TestErrorStyle(t *testing.T) {
	style := ErrorStyle()
	rendered := style.Render("Error message")

	if rendered == "" {
		t.Error("ErrorStyle produced empty render")
	}

	// Error style should contain the text even if formatting is not visible in tests
	if !strings.Contains(rendered, "Error message") {
		t.Error("ErrorStyle should contain the original text")
	}
}

func TestStyleConsistency(t *testing.T) {
	// Test that styles are consistent when called multiple times
	style1 := SuccessStyle()
	style2 := SuccessStyle()

	rendered1 := style1.Render("test")
	rendered2 := style2.Render("test")

	if rendered1 != rendered2 {
		t.Error("SuccessStyle should be consistent across calls")
	}

	// Same test for ErrorStyle
	errorStyle1 := ErrorStyle()
	errorStyle2 := ErrorStyle()

	errorRendered1 := errorStyle1.Render("test")
	errorRendered2 := errorStyle2.Render("test")

	if errorRendered1 != errorRendered2 {
		t.Error("ErrorStyle should be consistent across calls")
	}
}

func TestStylesWithEmptyString(t *testing.T) {
	// Test that styles handle empty strings gracefully
	styles := []struct {
		name  string
		style lipgloss.Style
	}{
		{"TitleStyle", TitleStyle},
		{"HelpStyle", HelpStyle},
		{"ItemStyle", ItemStyle},
	}

	for _, test := range styles {
		rendered := test.style.Render("")
		// Should not panic and should handle empty string
		_ = rendered
	}

	// Test function styles too
	successRendered := SuccessStyle().Render("")
	errorRendered := ErrorStyle().Render("")

	_ = successRendered
	_ = errorRendered
}

func TestStyleProperties(t *testing.T) {
	// Test that certain styles contain the expected text

	// TitleStyle should contain the text
	titleRendered := TitleStyle.Render("Title")
	if !strings.Contains(titleRendered, "Title") {
		t.Error("TitleStyle should contain the original text")
	}

	// SuccessStyle should contain the text
	successRendered := SuccessStyle().Render("Success")
	if !strings.Contains(successRendered, "Success") {
		t.Error("SuccessStyle should contain the original text")
	}

	// ErrorStyle should contain the text
	errorRendered := ErrorStyle().Render("Error")
	if !strings.Contains(errorRendered, "Error") {
		t.Error("ErrorStyle should contain the original text")
	}
}
