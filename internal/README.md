# internal

Caravan's implementation, split into two packages so the Cobra command
wiring stays separate from the actual cluster mechanics.

## cli

Cobra commands. Each file registers its commands on `rootCmd` via `init()`
and delegates to `internal/cluster` — no cluster logic lives here.

- **root.go** — `rootCmd` + `Execute()`, called from `main.go`.
- **cluster.go** — `cluster up|down|status`.
- **submit.go** — `submit <script.sh>`.

## cluster

The actual engine: picks a container runtime, extracts the embedded Slurm
scaffold, and drives it.

- **engine.go** — `Engine` interface (`Up`/`Down`/`Ps`/`Exec`/`Logs`) with
  `DockerEngine`/`PodmanEngine` implementations. `NewEngine()` autodetects
  Docker first, then Podman, overridable via `CARAVAN_ENGINE`/`CARAVAN_COMPOSE`.
- **extract.go** — `Scaffold` (wired in from `main.go`'s embedded
  `slurm-cluster/`) and `Extract()`, which writes it to `CARAVAN_DIR`
  (default `~/.caravan/cluster`).
- **compose.go** — compose file path + project name.
- **status.go** — `Up`/`Down`/`Status`, the top-level operations the CLI calls.
- **submit.go** — streams a script into `sbatch` on the controller.
- **util.go** — `exec.Command` wrappers shared by the engines.

See [slurm-cluster/README.md](../slurm-cluster/README.md) for the cluster
scaffold itself.
