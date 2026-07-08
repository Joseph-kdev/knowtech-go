package firebase

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveCredentialsSource(t *testing.T) {
	t.Run("uses explicit json from env", func(t *testing.T) {
		t.Setenv("GOOGLE_CREDS", "")
		t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", `{"type":"service_account"}`)

		path, data, err := resolveCredentialsSource()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if path != "" {
			t.Fatalf("expected no file path, got %q", path)
		}
		if !strings.Contains(string(data), `"type":"service_account"`) {
			t.Fatalf("expected json credentials data, got %q", string(data))
		}
	})

	t.Run("uses file path from env", func(t *testing.T) {
		dir := t.TempDir()
		file := filepath.Join(dir, "service.json")
		if err := os.WriteFile(file, []byte(`{"type":"service_account"}`), 0o600); err != nil {
			t.Fatalf("write file: %v", err)
		}

		t.Setenv("GOOGLE_CREDS", file)
		t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "")

		path, data, err := resolveCredentialsSource()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if path != file {
			t.Fatalf("expected file path %q, got %q", file, path)
		}
		if len(data) != 0 {
			t.Fatalf("expected no inline json data, got %q", string(data))
		}
	})
}
