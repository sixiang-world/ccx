package uipreferences_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BenedictKing/ccx/desktop/internal/uipreferences"
)

func TestLoadReturnsDefaultsWhenMissing(t *testing.T) {
	dir := t.TempDir()

	prefs, exists, err := uipreferences.Load(dir)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if exists {
		t.Fatalf("expected no existing preferences")
	}
	if prefs.Locale != "" {
		t.Fatalf("expected empty locale, got %s", prefs.Locale)
	}
	if prefs.Manual {
		t.Fatalf("expected manual false")
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()

	if err := uipreferences.Save(dir, uipreferences.Preferences{Locale: "zh", Manual: true}); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(dir, "ui-preferences.json"))
	if err != nil {
		t.Fatalf("read preferences file failed: %v", err)
	}
	if string(raw) == "" {
		t.Fatalf("preferences file is empty")
	}

	prefs, exists, err := uipreferences.Load(dir)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if !exists {
		t.Fatalf("expected preferences to exist")
	}
	if prefs.Locale != uipreferences.LocaleChinese {
		t.Fatalf("expected zh-CN, got %s", prefs.Locale)
	}
	if !prefs.Manual {
		t.Fatalf("expected manual true")
	}
}

func TestLoadIgnoresCorruptedFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "ui-preferences.json"), []byte("{invalid"), 0o644); err != nil {
		t.Fatalf("write corrupted file failed: %v", err)
	}

	prefs, exists, err := uipreferences.Load(dir)
	if err != nil {
		t.Fatalf("expected nil error for corrupted file, got %v", err)
	}
	if exists {
		t.Fatalf("expected no existing preferences for corrupted file")
	}
	if prefs.Locale != "" {
		t.Fatalf("expected empty locale for corrupted file, got %s", prefs.Locale)
	}
}

func TestNormalizeLocaleMapping(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"en", uipreferences.LocaleEnglish},
		{"en-US", uipreferences.LocaleEnglish},
		{"zh", uipreferences.LocaleChinese},
		{"zh_CN", uipreferences.LocaleChinese},
		{"zh-CN", uipreferences.LocaleChinese},
		{"zh-Hans", uipreferences.LocaleChinese},
		{"zh_Hans_CN", uipreferences.LocaleChinese},
		{"zh_CN.UTF-8", uipreferences.LocaleChinese},
		{"fr-FR", uipreferences.LocaleEnglish},
		{"  ", ""},
		{"", ""},
	}
	for _, tt := range tests {
		if got := uipreferences.NormalizeLocale(tt.input); got != tt.want {
			t.Errorf("NormalizeLocale(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
