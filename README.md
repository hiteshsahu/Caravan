# Caravan 🐪🐫🛕

> **A CLI for GPU Slurm stand up a cluster and run workloads on it.**

![](./img/cover.jpg)

### Purpose
Slurm is powerful but a chore to operate and submit to. 

Caravan is a single binary that carries a GPU Slurm cluster *inside it* and makes it a one-liner to bring up —
then (next) submit jobs, track experiments, and rerun workloads.

Its completing project to:
- [squint](https://github.com/hiteshsahu/squint) watching the queue and 
- [gpu-lens](https://github.com/hiteshsahu/gpu-lens) watching the GPUs.

> Caravan **uses** Slurm — it doesn't replace it. Slurm stays the scheduler;

> Caravan is the control plane and developer experience around it.

---

## Quick start

```bash
# Run the built-in CLI directly from source
go run . cluster up

go run . cluster status

go run . cluster down
```

## Cluster up

### Prerequisites

- **Go 1.22+**

- ### Container engine

   Caravan works with `Docker` or `Podman`, it auto-detects (`Docker` first, then
  `Podman`) and uses the matching Compose.

  On Podman it prefers `podman-compose` if
  installed, else `podman compose`. Override either:

 
![](./img/illustration.jpeg)


### 1. Using Released Binary

```bash
go install github.com/hiteshsahu/caravan@latest
```

Setup Slurm cluster:

```bash
# builds the image, starts controller + 2 GPU nodes
caravan cluster up

# container state, then `sinfo`
caravan cluster status

# stop the cluster
caravan cluster down
# stop and wipe state
caravan cluster down -v
```

### Using Podman

Caravan auto-detects `docker` first, then `podman` if available.

Force Podman explicitly:

```bash
CARAVAN_ENGINE=podman caravan cluster up
CARAVAN_COMPOSE="podman compose" caravan cluster up
```

On macOS make sure the Podman VM is running first:

```bash
podman machine start
```

---

### 2. Using Local Binary

Build the binary after changes in Go files:

```bash
go build -o caravan .
```

Start the cluster from the local binary:

```bash
./caravan cluster up
```

`caravan cluster up` writes an embedded Slurm scaffold to `~/.caravan/cluster` and runs `docker compose` against it.

You can override the scaffold location with `CARAVAN_DIR`:

```bash
CARAVAN_DIR=/tmp/caravan/cluster ./caravan cluster up
```

The two compute nodes advertise `gpu:4` each as **fake, count-only GPUs** — real GPU scheduling, no hardware needed (no `nvidia-smi` telemetry).

**Check status**

```bash
./caravan cluster status
```

**Result:**

```text

hitesh@Mac Caravan % ./caravan cluster status
CONTAINER ID  IMAGE                           COMMAND     CREATED         STATUS                 PORTS       NAMES
ce05d9f1bbd9  localhost/caravan-slurm:latest  slurmd      15 seconds ago  Up 1 second (healthy)              c1
0ef798e4b97c  localhost/caravan-slurm:latest  slurmd      15 seconds ago  Up 1 second (healthy)              c2
2f3e0982f69c  localhost/caravan-slurm:latest  slurmctld   14 seconds ago  Up 14 seconds                      slurmctld

PARTITION AVAIL  TIMELIMIT  NODES  STATE NODELIST
gpu*         up   infinite      2    unk c[1-2]
                            c2
```

**Clean up**

```bash  
  ./caravan cluster down       # add -v to also remove volumes  
```  

---


## 👨‍💻 DEVELOP

###  ⚙️ Install dependencies
```bash
    # Install dependencies
    go mod tidy
```

###  🧪 Build & Test

Tests are run as part of CI itself.

 ```bash  
    # Formatting go file
    gofmt -w . 
    
    # Linting
    go vet ./... 
    
    # recursively compiles all packages
    go build ./...   

```

### ▶️ Run

``` bash
    # Run the Engine
    go run .                 # mock
    
```


---

## 📁 Folder Structure

```
  caravan/
  ├── main.go
  ├── internal/
  │   ├── cli/                 # cobra commands
  │   │   ├── root.go
  │   │   └── cluster.go       # caravan cluster up|down|status
  │   └── cluster/
  │       ├── cluster.go       # extract scaffold + docker compose
  │       └── assets/          # the GPU Slurm cluster, embedded in the binary
  │           ├── Dockerfile · entrypoint.sh
  │           ├── docker-compose.yml
  │           └── slurm.conf · gres.conf
```

The cluster definition is embedded with `//go:embed`, so the binary is
self-contained — there's no separate cluster repo to clone. `caravan cluster up`
extracts it and runs it.

----

## 🗺️ Roadmap

Caravan grows from "runs a cluster" to "runs your work on it."

- **cluster** *(here)* — `up` / `down` / `status`. Docker today; cloud / bare-metal
  backends later, behind the same interface.
- **submit** — `caravan submit job.yaml` → generates an sbatch script → submits to
  the cluster → records the job. `status` / `logs` to follow it.
- **rerun** — re-launch a past job by id, reproducibly.
- **exp** — group runs into experiments and compare them.

Each new capability sits behind a `Backend` interface (Slurm today), so swapping
the execution target later doesn't touch the CLI above it.

---

## License
*© 2026 [Hitesh Kumar Sahu](https://hiteshsahu.com) · Licensed under [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0)*

