package configservice

import (
	"strings"
	"testing"
)

func TestComputeTextDiff_EmptyBoth(t *testing.T) {
	result := computeTextDiff("test.txt", "", "")
	if len(result.Lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(result.Lines))
	}
}

func TestComputeTextDiff_Create(t *testing.T) {
	result := computeTextDiff("test.txt", "", "hello\nworld\n")
	if result.Action != "create" {
		t.Fatalf("action = %q, want create", result.Action)
	}
	if len(result.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(result.Lines))
	}
	for _, l := range result.Lines {
		if l.Type != "added" {
			t.Errorf("line type = %q, want added", l.Type)
		}
	}
}

func TestComputeTextDiff_Delete(t *testing.T) {
	result := computeTextDiff("test.txt", "hello\n", "")
	if result.Action != "delete" {
		t.Fatalf("action = %q, want delete", result.Action)
	}
	if len(result.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(result.Lines))
	}
	if result.Lines[0].Type != "removed" {
		t.Errorf("line type = %q, want removed", result.Lines[0].Type)
	}
}

func TestComputeTextDiff_Modify(t *testing.T) {
	before := "line1\nline2\nline3\n"
	after := "line1\nmodified\nline3\n"
	result := computeTextDiff("test.txt", before, after)
	if result.Action != "modify" {
		t.Fatalf("action = %q, want modify", result.Action)
	}
	var context, added, removed int
	for _, l := range result.Lines {
		switch l.Type {
		case "context":
			context++
		case "added":
			added++
		case "removed":
			removed++
		}
	}
	if context != 2 {
		t.Errorf("context lines = %d, want 2", context)
	}
	if added != 1 {
		t.Errorf("added lines = %d, want 1", added)
	}
	if removed != 1 {
		t.Errorf("removed lines = %d, want 1", removed)
	}
}

func TestComputeJSONDiff_Create(t *testing.T) {
	after := map[string]any{"key": "value"}
	result := computeJSONDiff("config.json", nil, after)
	if result.Action != "create" {
		t.Fatalf("action = %q, want create", result.Action)
	}
	found := false
	for _, l := range result.Lines {
		if l.Type == "added" && strings.Contains(l.Content, `"key"`) {
			found = true
		}
	}
	if !found {
		t.Error("expected added line containing 'key'")
	}
}

func TestComputeJSONDiff_ModifyEnv(t *testing.T) {
	before := map[string]any{
		"env": map[string]any{
			"ANTHROPIC_BASE_URL": "http://old:3000",
		},
	}
	after := map[string]any{
		"env": map[string]any{
			"ANTHROPIC_BASE_URL": "http://new:3000",
		},
	}
	result := computeJSONDiff("settings.json", before, after)
	if result.Action != "modify" {
		t.Fatalf("action = %q, want modify", result.Action)
	}
	hasOld := false
	hasNew := false
	for _, l := range result.Lines {
		if l.Type == "removed" && strings.Contains(l.Content, "old:3000") {
			hasOld = true
		}
		if l.Type == "added" && strings.Contains(l.Content, "new:3000") {
			hasNew = true
		}
	}
	if !hasOld || !hasNew {
		t.Errorf("expected removed old url and added new url; hasOld=%v hasNew=%v", hasOld, hasNew)
	}
}

func TestMaskSensitiveValue(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"空字符串", "", ""},
		{"短值", "abc", "***"},
		{"正好11字符", "12345678901", "***"},
		{"12字符", "123456789012", "123***9012"},
		{"长值", "sk-abcdef12345678", "sk-***5678"},
		{"带空格", "  sk-test12345678  ", "sk-***5678"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := maskSensitiveValue(c.input)
			if got != c.want {
				t.Errorf("maskSensitiveValue(%q) = %q, want %q", c.input, got, c.want)
			}
		})
	}
}

func TestMaskMapSensitiveKeys(t *testing.T) {
	data := map[string]any{
		"ANTHROPIC_API_KEY":  "sk-very-long-key-12345678",
		"ANTHROPIC_BASE_URL": "http://localhost:3000",
		"OTHER_SECRET":       "should-not-be-masked",
	}
	result := maskMapSensitiveKeys(data, "ANTHROPIC_API_KEY")
	if result["ANTHROPIC_BASE_URL"] != "http://localhost:3000" {
		t.Error("non-sensitive field should not be modified")
	}
	if result["OTHER_SECRET"] != "should-not-be-masked" {
		t.Error("unspecified field should not be masked")
	}
	if result["ANTHROPIC_API_KEY"] == "sk-very-long-key-12345678" {
		t.Error("ANTHROPIC_API_KEY should be masked")
	}
	masked := result["ANTHROPIC_API_KEY"].(string)
	if !strings.HasPrefix(masked, "sk-") {
		t.Errorf("masked value should start with prefix, got %q", masked)
	}
	// 原始 map 不应被修改
	if data["ANTHROPIC_API_KEY"] != "sk-very-long-key-12345678" {
		t.Error("original map should not be modified")
	}
}

func TestMaskJSONSensitiveKeys(t *testing.T) {
	data := map[string]any{
		"env": map[string]any{
			"ANTHROPIC_BASE_URL": "http://localhost:3000",
			"ANTHROPIC_API_KEY":  "sk-abcdef123456789",
		},
		"permissions": map[string]any{},
	}
	result := maskJSONSensitiveKeys(data)
	env := result["env"].(map[string]any)
	if env["ANTHROPIC_BASE_URL"] != "http://localhost:3000" {
		t.Error("non-sensitive env field should not be masked")
	}
	if env["ANTHROPIC_API_KEY"] == "sk-abcdef123456789" {
		t.Error("ANTHROPIC_API_KEY should be masked")
	}
}

func TestMaskTextSensitiveValues(t *testing.T) {
	content := `model_provider = "ccx"
OPENAI_API_KEY = "sk-test1234567890"`
	keyValues := map[string]string{
		"OPENAI_API_KEY": "sk-test1234567890",
	}
	result := maskTextSensitiveValues(content, keyValues)
	if strings.Contains(result, "sk-test1234567890") {
		t.Error("original key value should be masked in output")
	}
	if !strings.Contains(result, "sk-***7890") {
		t.Errorf("expected masked value in output, got:\n%s", result)
	}
	if !strings.Contains(result, `model_provider = "ccx"`) {
		t.Error("non-sensitive content should be preserved")
	}
}

func TestLcsDiff_IdenticalContent(t *testing.T) {
	lines := lcsDiff([]string{"a", "b", "c"}, []string{"a", "b", "c"})
	for _, l := range lines {
		if l.Type != "context" {
			t.Errorf("type = %q, want context", l.Type)
		}
	}
}

func TestLcsDiff_AllAdded(t *testing.T) {
	lines := lcsDiff(nil, []string{"x", "y"})
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	for _, l := range lines {
		if l.Type != "added" {
			t.Errorf("type = %q, want added", l.Type)
		}
	}
}

func TestLcsDiff_AllRemoved(t *testing.T) {
	lines := lcsDiff([]string{"x", "y"}, nil)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	for _, l := range lines {
		if l.Type != "removed" {
			t.Errorf("type = %q, want removed", l.Type)
		}
	}
}

func TestLcsDiff_NoCommonLines(t *testing.T) {
	lines := lcsDiff([]string{"a", "b"}, []string{"x", "y"})
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}
	removed, added := 0, 0
	for _, l := range lines {
		switch l.Type {
		case "removed":
			removed++
		case "added":
			added++
		default:
			t.Errorf("unexpected type %q", l.Type)
		}
	}
	if removed != 2 || added != 2 {
		t.Errorf("removed=%d added=%d, want 2/2", removed, added)
	}
}

func TestLcsDiff_SingleLineReplace(t *testing.T) {
	lines := lcsDiff([]string{"old"}, []string{"new"})
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0].Type != "removed" || lines[0].Content != "old" {
		t.Errorf("line 0: type=%q content=%q, want removed/old", lines[0].Type, lines[0].Content)
	}
	if lines[1].Type != "added" || lines[1].Content != "new" {
		t.Errorf("line 1: type=%q content=%q, want added/new", lines[1].Type, lines[1].Content)
	}
}

func TestLcsDiff_MiddleInsert(t *testing.T) {
	lines := lcsDiff([]string{"a", "c"}, []string{"a", "b", "c"})
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0].Type != "context" || lines[0].Content != "a" {
		t.Errorf("line 0: want context/a, got %q/%q", lines[0].Type, lines[0].Content)
	}
	if lines[1].Type != "added" || lines[1].Content != "b" {
		t.Errorf("line 1: want added/b, got %q/%q", lines[1].Type, lines[1].Content)
	}
	if lines[2].Type != "context" || lines[2].Content != "c" {
		t.Errorf("line 2: want context/c, got %q/%q", lines[2].Type, lines[2].Content)
	}
}

func TestComputeTextDiff_IdenticalContent(t *testing.T) {
	before := "line1\nline2\n"
	after := "line1\nline2\n"
	result := computeTextDiff("test.txt", before, after)
	if result.Action != "modify" {
		t.Fatalf("action = %q, want modify", result.Action)
	}
	for _, l := range result.Lines {
		if l.Type != "context" {
			t.Errorf("type = %q, want context for identical content", l.Type)
		}
	}
}
