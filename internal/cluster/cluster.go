package cluster

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// The GPU Slurm cluster scaffold travels inside the binary.
//
//go:embed assets/*
var assets embed.FS

const project = "caravan"

// scaffoldDir is where the embedded assets get written before compose runs.
// Override with CARAVAN_DIR.
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

func composeFile(dir string) string { return filepath.Join(dir, "docker-compose.yml") }

func hasBin(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// resolveEngine picks a container engine and a compose command. Defaults:
// Docker if present, else Podman. Override with CARAVAN_ENGINE (docker|podman)
// and CARAVAN_COMPOSE (e.g. "podman-compose" or "docker compose").
func resolveEngine() (cli string, compose []string, err error) {
	cli = os.Getenv("CARAVAN_ENGINE")
	if cli == "" {
		switch {
		case hasBin("docker"):
			cli = "docker"
		case hasBin("podman"):
			cli = "podman"
		case hasBin("podman-compose"):
			cli = "podman"
		default:
			return "", nil, fmt.Errorf("no container engine on PATH — install Docker or Podman, or set CARAVAN_ENGINE")
		}
	}

	switch {
	case os.Getenv("CARAVAN_COMPOSE") != "":
		compose = strings.Fields(os.Getenv("CARAVAN_COMPOSE"))
	case cli == "podman" && hasBin("podman-compose"):
		compose = []string{"podman-compose"} // self-contained; avoids needing a compose provider
	default:
		compose = []string{cli, "compose"}
	}
	return cli, compose, nil
}

// Extract writes the embedded scaffold into dir.
func Extract(dir string) error {
	return fs.WalkDir(assets, "assets", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == "assets" {
			return nil
		}
		rel := strings.TrimPrefix(p, "assets/")
		target := filepath.Join(dir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := assets.ReadFile(p)
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

// composeArgv builds the full argv for a compose invocation.
func composeArgv(compose []string, dir string, extra ...string) (string, []string) {
	args := append([]string{}, compose[1:]...)
	args = append(args, "-p", project, "-f", composeFile(dir))
	args = append(args, extra...)
	return compose[0], args
}

// Up extracts the scaffold and brings the cluster online.
func Up() error {
	dir, err := scaffoldDir()
	if err != nil {
		return err
	}
	cli, compose, err := resolveEngine()
	if err != nil {
		return err
	}
	fmt.Printf("→ engine: %s (compose: %s)\n", cli, strings.Join(compose, " "))
	fmt.Printf("→ scaffolding GPU Slurm cluster in %s\n", dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if err := Extract(dir); err != nil {
		return fmt.Errorf("extract scaffold: %w", err)
	}
	fmt.Println("→ starting cluster (first run builds the image)…")
	name, args := composeArgv(compose, dir, "up", "-d", "--build")
	return runIn(dir, name, args...)
}

// Down stops the cluster; volumes also wipes its state.
func Down(volumes bool) error {
	dir, err := scaffoldDir()
	if err != nil {
		return err
	}
	_, compose, err := resolveEngine()
	if err != nil {
		return err
	}
	extra := []string{"down"}
	if volumes {
		extra = append(extra, "--volumes")
	}
	name, args := composeArgv(compose, dir, extra...)
	return runIn(dir, name, args...)
}

// Status prints container state, then Slurm node state.
func Status() error {
	dir, err := scaffoldDir()
	if err != nil {
		return err
	}
	cli, compose, err := resolveEngine()
	if err != nil {
		return err
	}
	name, args := composeArgv(compose, dir, "ps")
	_ = runIn(dir, name, args...)
	fmt.Println()
	return run(cli, "exec", "slurmctld", "sinfo")
}

func Submit(scriptPath string) error {
	dir, err := scaffoldDir()
	if err != nil {
		return err
	}
	_, compose, err := resolveEngine()
	if err != nil {
		return err
	}
	if _, err := os.Stat(scriptPath); err != nil {
		return err
	}
	scriptFile, err := os.Open(scriptPath)
	if err != nil {
		return err
	}
	defer scriptFile.Close()

	fmt.Printf("→ submitting %s to local Slurm cluster in %s\n", scriptPath, dir)
	name, args := composeArgv(compose, dir, "exec", "slurmctld", "sh", "-c", "cat > /tmp/caravan-job.sh && sbatch --parsable /tmp/caravan-job.sh")
	return runInWithStdin(dir, name, args, scriptFile)
}

func run(name string, args ...string) error { return runIn("", name, args...) }

func runIn(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runInWithStdin(dir, name string, args []string, stdin io.Reader) error {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = stdin
	return cmd.Run()
}
