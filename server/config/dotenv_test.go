package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFirstDotEnvPrefersFirstExistingCandidate(t *testing.T) {
	dir := t.TempDir()
	parent := filepath.Join(dir, ".env")
	if err := os.WriteFile(parent, []byte("CRUX_TEST_DOTENV=from-parent\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Unsetenv("CRUX_TEST_DOTENV") })

	missing := filepath.Join(dir, "missing", ".env")
	got, err := loadFirstDotEnv(missing, parent)
	if err != nil {
		t.Fatalf("loadFirstDotEnv returned error: %v", err)
	}
	if got != parent {
		t.Fatalf("loaded %q, want %q", got, parent)
	}
	if v := os.Getenv("CRUX_TEST_DOTENV"); v != "from-parent" {
		t.Fatalf("CRUX_TEST_DOTENV = %q, want from-parent", v)
	}
}

func TestLoadFirstDotEnvNoFileIsNotAnError(t *testing.T) {
	got, err := loadFirstDotEnv(filepath.Join(t.TempDir(), ".env"))
	if err != nil || got != "" {
		t.Fatalf("got (%q, %v), want (\"\", nil)", got, err)
	}
}

func TestLoadDotEnvFromServerDirUsesRootFile(t *testing.T) {
	root := t.TempDir()
	serverDir := filepath.Join(root, "server")
	if err := os.Mkdir(serverDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte("CRUX_TEST_DOTENV_ROOT=yes\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Unsetenv("CRUX_TEST_DOTENV_ROOT") })
	t.Chdir(serverDir)

	if err := LoadDotEnv(); err != nil {
		t.Fatalf("LoadDotEnv: %v", err)
	}
	if v := os.Getenv("CRUX_TEST_DOTENV_ROOT"); v != "yes" {
		t.Fatalf("CRUX_TEST_DOTENV_ROOT = %q, want yes", v)
	}
}
