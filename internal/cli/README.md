# cli

Cobra command definitions. Each file registers its command(s) on `rootCmd`
via `init()` and delegates straight to `internal/cluster` — no cluster logic
lives here, just argument/flag parsing and wiring.

- **root.go** — `rootCmd` (the `caravan` command) and `Execute()`, called
  from `main.go`.
- **cluster.go** — `cluster up|down|status`. `down` takes `-v`/`--volumes`
  to also wipe cluster state.
- **submit.go** — `submit <script.sh>`, takes exactly one argument.

## Commands

| Command                  | Flags             | Description                                  |
|---------------------------|-------------------|-----------------------------------------------|
| `caravan cluster up`      | —                 | Build + start the cluster (controller + 2 nodes) |
| `caravan cluster down`    | `-v, --volumes`   | Stop the cluster, optionally wiping volumes  |
| `caravan cluster status`  | —                 | Show container state + `sinfo`               |
| `caravan submit <script>` | —                 | Stream a script into `sbatch` on the controller |

To add a command: create a file here with a `cobra.Command` + an `init()`
that calls `rootCmd.AddCommand(...)` (or a parent command's `AddCommand`,
as `cluster.go` does for its subcommands), and have `RunE` call into
`internal/cluster`.
