package desktoplauncher

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunLegacyMigratorSkipsCompatSymlink(t *testing.T) {
	root := t.TempDir()
	launcher := filepath.Join(root, "reasonix-launcher")
	if err := os.WriteFile(launcher, []byte("launcher"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("reasonix-launcher", filepath.Join(root, "reasonix-guard")); err != nil {
		t.Fatal(err)
	}
	if err := runLegacyMigratorIfNeeded(root); err != nil {
		t.Fatalf("compat symlink must be ignored, got %v", err)
	}
}
