package uipreferences

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const fileName = "ui-preferences.json"

const (
	LocaleEnglish = "en"
	LocaleChinese = "zh-CN"
)

type Preferences struct {
	Locale string `json:"locale,omitempty"`
	Manual bool   `json:"manual,omitempty"`
}

func Load(dataDir string) (Preferences, bool, error) {
	path := filepath.Join(dataDir, fileName)
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Preferences{}, false, nil
		}
		return Preferences{}, false, err
	}
	var prefs Preferences
	if err := json.Unmarshal(raw, &prefs); err != nil {
		return Preferences{}, false, nil
	}
	prefs.Locale = NormalizeLocale(prefs.Locale)
	if prefs.Locale == "" {
		return Preferences{}, false, nil
	}
	return prefs, true, nil
}

func Save(dataDir string, prefs Preferences) error {
	prefs.Locale = NormalizeLocale(prefs.Locale)
	if prefs.Locale == "" {
		prefs.Locale = LocaleEnglish
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dataDir, "ui-preferences-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()
	if err := json.NewEncoder(tmp).Encode(prefs); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	target := filepath.Join(dataDir, fileName)
	return os.Rename(tmpName, target)
}

func NormalizeLocale(raw string) string {
	locale := strings.TrimSpace(raw)
	if locale == "" {
		return ""
	}
	locale = strings.ReplaceAll(locale, "_", "-")
	if dot := strings.Index(locale, "."); dot >= 0 {
		locale = locale[:dot]
	}
	lower := strings.ToLower(locale)
	if lower == "en" || strings.HasPrefix(lower, "en-") {
		return LocaleEnglish
	}
	if lower == "zh" || strings.HasPrefix(lower, "zh-") {
		return LocaleChinese
	}
	return LocaleEnglish
}
