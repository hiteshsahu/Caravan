package cluster

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Scaffold is the embedded GPU Slurm cluster scaffold. go:embed directives
// can't reach outside their own file's directory tree, so slurm-cluster/
// lives at the repo root and main wires it in here rather than this
// package embedding it directly.
var Scaffold fs.FS

// scaffoldDir is where the embedded scaffold gets written before compose
// runs. Override with CARAVAN_DIR.
func scaffoldDir() (string, error) {
	if d := os.Getenv("CARAVAN_DIR"); d != "" {
		return d, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".caravan", "cluster"), nil
}

// Extract writes the embedded scaffold into dir.
func Extract(dir string) error {
	return fs.WalkDir(Scaffold, "slurm-cluster", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == "slurm-cluster" {
			return nil
		}
		rel := strings.TrimPrefix(p, "slurm-cluster/")
		target := filepath.Join(dir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := fs.ReadFile(Scaffold, p)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		mode := os.FileMode(0o644)
		if strings.HasSuffix(rel, ".sh") {
			mode = 0o755
		}
		return os.WriteFile(target, b, mode)
	})
}
