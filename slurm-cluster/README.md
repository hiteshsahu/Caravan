# Slurm Cluster Scaffold

This is the GPU Slurm cluster that Caravan carries inside its binary. It's
embedded via `//go:embed slurm-cluster/*` in [embed.go](../embed.go) at the
repo root, then extracted to `CARAVAN_DIR` (default `~/.caravan/cluster`) and
run with `docker`/`podman compose` by `caravan cluster up`.

Editing files here changes what gets embedded in the next build — there's no
separate generated copy to keep in sync.

## Files

- **Dockerfile** — Ubuntu 24.04 image with `slurm-wlm` + `munge` installed.
- **entrypoint.sh** — starts `munged`, then execs `slurmctld` or `slurmd`
  depending on the container's `command`.
- **docker-compose.yml** — one `slurmctld` controller plus two `slurmd`
  compute nodes (`c1`, `c2`), each advertising `gpu:4` as fake, count-only
  GPUs (no real hardware, no `nvidia-smi`) by default.
- **slurm.conf** — cluster config. `NodeName` CPUs/RealMemory must not
  exceed what Docker actually gives the container, or the node registers
  invalid/drained.
- **cgroup.conf** — `CgroupPlugin=cgroup/v2` + `IgnoreSystemd=yes`, since
  these containers have no systemd/dbus to manage cgroup scopes for them.
- **gres.conf** — declares the fake `gpu:4` GRES per node.
- **docker-compose.gpu.yml**, **slurm.gpu.conf**, **gres.gpu.conf** — the
  `CARAVAN_GPU=real` overlay (see below); inert unless that env var is set.

## Known quirks

- `c1`/`c2` run `privileged: true` + `cgroup: host` + `cgroup_parent: "/"`.
  slurmd's `cgroup/v2` plugin assumes its own container cgroup sits directly
  under the cgroup mount root (so it can create a `system.slice/<scope>`
  sibling) — true by default on hosts using the `systemd` cgroup driver, but
  not on hosts using `cgroupfs` (e.g. Docker Desktop on Windows/WSL2). The
  compose settings force that layout regardless of host driver.
- The Dockerfile wraps the real `slurmstepd` binary in `unshare --cgroup`
  (see its comment). Without this, slurmstepd inherits slurmd's
  already-relocated cgroup and reapplies the same relocation logic relative
  to it, doubling the path and failing every job launch with
  `_forkexec_slurmstepd: ... Resource temporarily unavailable`. This bites
  regardless of GPU mode — fake mode just never reached it, since its nodes
  are drained before any job gets that far (see below).
- Fake-GPU nodes report `gres/gpu count reported lower than configured
  (0 < 4)` and drain a few seconds after slurmd starts, since these are
  file-less fake GPUs. Submitted jobs will queue but stay pending — not yet
  fixed.

## Real GPU (opt-in): `CARAVAN_GPU=real`

```bash
CARAVAN_GPU=real ./caravan cluster up
CARAVAN_GPU=real ./caravan submit workloads/gpu_example.sh
```

Layers `docker-compose.gpu.yml` on top of the default compose file (same
`-p`/project, just additional `-f`): `c1` gets the host's real GPU via
compose's `deploy.resources.reservations.devices` (the same mechanism as
`docker run --gpus`), and `slurmctld`/`c1`/`c2` mount `slurm.gpu.conf` +
`gres.gpu.conf` instead of the defaults. `c2` keeps the fake `gpu:4` —
there's only one physical GPU, so only `c1` gets it.

Why `gpu/generic` with `File=/dev/dxg` instead of `AutoDetect=nvml`:
Ubuntu's `slurm-wlm-basic-plugins` package ships `gpu_generic.so` and
`gpu_nrt.so`, but no `gpu_nvml.so` — NVML autodetect isn't available at
all, regardless of whether the library is present. Separately, GPU
passthrough into WSL2 containers exposes no `/dev/nvidia*` device nodes,
only `/dev/dxg` — so `gpu/generic` with that device file is what's actually
available here, on any host. No CUDA base image is needed either: Docker's
NVIDIA Container Toolkit injects `nvidia-smi`/`libnvidia-ml.so` into
whatever image you run once the device is reserved.

Requires an NVIDIA GPU with Docker Desktop GPU support enabled (verify with
`docker run --rm --gpus all <any-image> nvidia-smi` first).
