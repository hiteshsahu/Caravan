# cluster

The engine behind the CLI: picks a container runtime, extracts the embedded
Slurm scaffold, and drives it via compose.

- **engine.go** — `Engine` interface (`Up`/`Down`/`Ps`/`Exec`/`Logs`).
  `DockerEngine` and `PodmanEngine` both embed `baseEngine`, which differs
  only in which binary and compose command it uses. `NewEngine()` picks
  Docker if present, else Podman; override with `CARAVAN_ENGINE`
  (`docker`/`podman`) and `CARAVAN_COMPOSE` (e.g. `"podman-compose"` or
  `"docker compose"`).
- **extract.go** — `Scaffold` holds the embedded `slurm-cluster/` tree,
  wired in from `main.go` (a `go:embed` directive can't reach outside its
  own file's directory, so this package can't embed it directly).
  `Extract()` writes it to `dir` (preserving `+x` on `.sh` files).
- **compose.go** — `composeFile(dir)`/`composeGPUFile(dir)`, the fixed
  compose project name (`caravan`, used for `-p` so multiple checkouts
  don't collide), and `GPUEnabled()` (`CARAVAN_GPU=real`), which
  `engine.go`'s `composeArgv` checks to layer `docker-compose.gpu.yml` on
  top of the default compose file.
- **status.go** — `Up`/`Down`/`Status`: resolve `scaffoldDir()`
  (`CARAVAN_DIR`, default `~/.caravan/cluster`), get an `Engine`, and call
  through. These are what the `cli` package's `RunE` functions call.
- **submit.go** — `Submit(scriptPath)` streams a script's contents into
  `sh -c "cat > /tmp/caravan-job.sh && sbatch --parsable ..."` on the
  `slurmctld` container via `Engine.Exec`.
- **util.go** — thin `exec.Command` wrappers (`run`/`runIn`/
  `runInWithStdin`) shared by the engines; stdout/stderr always stream to
  the calling terminal.

See [../../slurm-cluster/README.md](../../slurm-cluster/README.md) for the
scaffold these extract and run.
