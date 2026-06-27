package cluster

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Engine drives a container runtime (Docker or Podman): it brings the
// compose stack up/down, reports container state, and runs commands inside
// a named container.
type Engine interface {
	Name() string
	Up(dir string) error
	Down(dir string, volumes bool) error
	Ps(dir string) error
	Exec(container string, stdin io.Reader, args ...string) error
	Logs(container string, args ...string) error
}

func hasBin(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// baseEngine implements Engine. Docker and Podman only differ in which
// binary and compose command they use, both fixed at construction.
type baseEngine struct {
	name    string
	cli     string   // binary used for direct exec/logs, e.g. "docker"
	compose []string // argv prefix for compose, e.g. ["docker", "compose"]
}

func (e *baseEngine) Name() string { return e.name }

func (e *baseEngine) composeArgv(dir string, extra ...string) (string, []string) {
	args := append([]string{}, e.compose[1:]...)
	args = append(args, "-p", project, "-f", composeFile(dir))
	if GPUEnabled() {
		args = append(args, "-f", composeGPUFile(dir))
	}
	args = append(args, extra...)
	return e.compose[0], args
}

func (e *baseEngine) Up(dir string) error {
	name, args := e.composeArgv(dir, "up", "-d", "--build")
	return runIn(dir, name, args...)
}

func (e *baseEngine) Down(dir string, volumes bool) error {
	extra := []string{"down"}
	if volumes {
		extra = append(extra, "--volumes")
	}
	name, args := e.composeArgv(dir, extra...)
	return runIn(dir, name, args...)
}

func (e *baseEngine) Ps(dir string) error {
	name, args := e.composeArgv(dir, "ps")
	return runIn(dir, name, args...)
}

func (e *baseEngine) Exec(container string, stdin io.Reader, args ...string) error {
	if stdin == nil {
		stdin = os.Stdin
	}
	// -i forwards stdin to the exec'd process; without it docker/podman
	// close stdin immediately, which truncates anything piped in (e.g. Submit).
	full := append([]string{"exec", "-i", container}, args...)
	return runInWithStdin("", e.cli, full, stdin)
}

func (e *baseEngine) Logs(container string, args ...string) error {
	full := append([]string{"logs", container}, args...)
	return run(e.cli, full...)
}

// DockerEngine drives Docker + `docker compose`.
type DockerEngine struct{ baseEngine }

// PodmanEngine drives Podman, preferring the standalone podman-compose
// binary when present (self-contained; avoids needing a compose provider).
type PodmanEngine struct{ baseEngine }

// NewEngine picks a container engine and its compose command. Defaults:
// Docker if present, else Podman. Override with CARAVAN_ENGINE (docker|podman)
// and CARAVAN_COMPOSE (e.g. "podman-compose" or "docker compose").
func NewEngine() (Engine, error) {
	cli := os.Getenv("CARAVAN_ENGINE")
	if cli == "" {
		switch {
		case hasBin("docker"):
			cli = "docker"
		case hasBin("podman"), hasBin("podman-compose"):
			cli = "podman"
		default:
			return nil, fmt.Errorf("no container engine on PATH — install Docker or Podman, or set CARAVAN_ENGINE")
		}
	}

	var compose []string
	switch {
	case os.Getenv("CARAVAN_COMPOSE") != "":
		compose = strings.Fields(os.Getenv("CARAVAN_COMPOSE"))
	case cli == "podman" && hasBin("podman-compose"):
		compose = []string{"podman-compose"}
	default:
		compose = []string{cli, "compose"}
	}

	switch cli {
	case "docker":
		return &DockerEngine{baseEngine{name: cli, cli: cli, compose: compose}}, nil
	case "podman":
		return &PodmanEngine{baseEngine{name: cli, cli: cli, compose: compose}}, nil
	default:
		return nil, fmt.Errorf("unknown engine %q (want docker or podman)", cli)
	}
}
