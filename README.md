# Caravan 🐪🐫🛕

> ### Make it easy to spin SLURM Cluster and submit your HPC workloads to it.

![Caravan COver](./img/cover.jpg)

![Build](https://github.com/hiteshsahu/caravan/actions/workflows/ci.yaml/badge.svg)
![Release](https://img.shields.io/github/v/release/hiteshsahu/caravan)
![Go Version](https://img.shields.io/github/go-mod/go-version/hiteshsahu/caravan) ![GoReleaser](https://img.shields.io/badge/Built%20with-GoReleaser-blue)

![License](https://img.shields.io/github/license/hiteshsahu/caravan)
![Downloads](https://img.shields.io/github/downloads/hiteshsahu/caravan/total)
![GitHub stars](https://img.shields.io/github/stars/hiteshsahu/caravan?style=social)

## Why Caravan

Developing Slurm workloads usually requires access to an HPC cluster. Slurm is powerful but a chore to operate and submit to.

Caravan makes it simple by bundling a complete Slurm cluster into a single CLI so you can develop, test and debug jobs locally.

It uses Docker or Podman behind the scenes and works with fake GPUs, making it ideal for CI, workshops and local development.

The cluster definition is embedded with `//go:embed`, so the binary is self-contained — there's no separate cluster repo to clone. `caravan cluster up` extracts it and runs it via a `docker`/`podman` `Engine`.

> Caravan **uses** Slurm — it doesn't replace it. Slurm stays the scheduler;

> Caravan is the control plane and easier developer experience around it.

```mermaid
flowchart LR
  subgraph CaravanCLI["🐫 Caravan (CLI)"]
    CLI["caravan cluster / submit / status"]
  end

  CLI -->|"sbatch"| CTL
  CLI -->|"squeue · scontrol · sinfo"| CTL

  subgraph Slurm["Slurm cluster 🥤 <br/>— embedded, no accounting DB"]
  
    CTL["slurmctld 🧠<br/>controller & scheduler"]
    CTL --> N1["slurmd 🖥<br/>compute node"]
    CTL --> N2["slurmd 🖥ִׄ<br/>compute node"]
    N1 --> S1["slurmstepd → GPU 🧮"]
    N2 --> S2["slurmstepd → GPU 🧮"]
  end

  N1 -. "DCGM / nvidia-smi" .-> OBS["squint 🦝 · gpu-lens ֎"]
  N2 -. "DCGM / nvidia-smi" .-> OBS

  classDef ctl fill:#EEEDFE,stroke:#534AB7,color:#26215C;
  classDef compute fill:#E1F5EE,stroke:#0F6E56,color:#04342C;
  classDef obs fill:#F1EFE8,stroke:#5F5E5A,color:#2C2C2A,stroke-dasharray:5 4;
  class CTL ctl;
  class N1,N2,S1,S2 compute;
  class OBS obs;
```

Its completing project to:
- [squint](https://github.com/hiteshsahu/squint): TUI Dashboard to check workload & squatting GPUs
- [gpu-lens](https://github.com/hiteshsahu/gpu-lens) : Drop-in GPU + scheduler observability for clusters(SLURM+K8)


---


## ⚡ Prerequisites

📟 Requires **Go 1.22+**

  > choco install golang

  More detail on this [Medium Post](https://medium.com/@hiteshkrsahu/installing-go-on-windows-the-5-minute-guide-and-the-gotchas-nobody-mentions-878eb3ea2277)
  
  **Slurm** will be installed as Container Image

### 🐳 Container engine
Caravan can work with both `Docker` or `Podman`, it auto-detects (`Docker` first, then`Podman`) and uses the matching Compose.

On Podman it uses `podman-compose` if installed, otherwise it will try to use `podman compose`.

You can Force Podman explicitly:

```bash
CARAVAN_ENGINE=podman caravan cluster up
CARAVAN_COMPOSE="podman compose" caravan cluster up
```

On macOS make sure the Podman VM is running first:

```bash
podman machine start
```

##  OS Support

Recommended Platform

| Environment       | Recommendation                                              |
|-------------------|-------------------------------------------------------------|
| 🐧 Linux          | ⭐⭐⭐⭐⭐ Best for full development, testing, and production    |
| 🪟 Windows + WSL2 | ⭐⭐⭐⭐⭐ Best Windows experience, close to Linux               |
| 🍎 macOS          | ⭐⭐⭐⭐ Great for UI and mock-mode development                 |
| 🪟 Native Windows | ⭐⭐⭐ Good for CLI/UI development; use WSL2 for Linux tooling |


### Slurm + Caravan OS Compatibility

| Feature                                    | Linux | Windows + WSL2 |      Windows      |       macOS       |
|--------------------------------------------|:-----:|:--------------:|:-----------------:|:-----------------:|
| Build                                      |   ✅   |       ✅        |         ✅         |         ✅         |
| Run mock mode                              |   ✅   |       ✅        |         ✅         |         ✅         |
| TUI development                            |   ✅   |       ✅        |         ✅         |         ✅         |
| Unit tests                                 |   ✅   |       ✅        |         ✅         |         ✅         |
| Integration tests (mock)                   |   ✅   |       ✅        |         ✅         |         ✅         |
| Live Slurm (`squeue`, `sacct`, `scontrol`) |   ✅   |       ✅        | ⚠️ Remote cluster | ⚠️ Remote cluster |
| DCGM GPU metrics                           |   ✅   |       ✅        |         ❌         |         ❌         |
| `nvidia-smi` GPU metrics                   |   ✅   |       ✅        |         ✅         |    ⚠️ Limited     |
| Full end-to-end testing                    |   ✅   |       ⚠️       |         ❌         |         ❌         |


Notes
-  Live Slurm support in WSL2 works if you have access to a remote Linux Slurm cluster (via SSH) or a local Slurm installation inside WSL2.
-  DCGM requires machine with NVIDIA GPU with CUDA support and the NVIDIA WSL driver stack.
-  nvidia-smi is only available on Windows and inside WSL2 when using the NVIDIA WSL GPU drivers.

---

## Start Your Caravan 🐪🐫

**Once you have CLI ready, you can start cluster and submit jobs**

![Caravan Batch](./img/illustration.jpeg)

###  ✨ 1. Using pre compile released Binary (for production)

Quick start with released CLI binary

**On macOS**

```bash
  go install github.com/hiteshsahu/caravan@latest
  caravan
  
  caravan cluster up        # build + start (controller + 2 fake-GPU nodes)
  caravan cluster down      # stop (-v to also wipe volumes)
  caravan cluster status    # container state + sinfo
  caravan submit <script>   # stream a script into sbatch on the controller
```

### 

### ⚙️ 2. Build Locally (for devs)

Quick start with local binary. 

Note: Replace `caravan` with `./caravan` and you can use them for local CLI

📦 **On macOS**

```bash
  go build -o caravan .
  ./caravan
  
  ./caravan cluster up        # build + start (controller + 2 fake-GPU nodes)
  ./caravan cluster down      # stop (-v to also wipe volumes)
  ./caravan cluster status    # container state + sinfo
  ./caravan submit <script>   # stream a script into sbatch on the controller
```


⊞ **On Windows**

**Note:** for PowerShell users: use `./caravan.exe` instead of `./caravan`.

```bash
   go build -o caravan.exe .
   ./caravan.exe --help
  
  ./caravan.exe cluster up        # build + start (controller + 2 fake-GPU nodes)
  ./caravan.exe cluster status    # container state + sinfo
  ./caravan.exe submit workloads/submit_example.sh
  ./caravan.exe cluster down      # stop (-v to also wipe volumes)
```

---


## How it works

### ▶️ Starting Slurm Cluster 

Writes an embedded Slurm scaffold to 
- MacOS:   `~/.caravan/cluster` 
- Windows: `%USERPROFILE%\.caravan\cluster`

and runs `docker`/`podman compose` against it.

```bash
  ./caravan cluster up
```

You can override scafold to your desired directory by passing `CARAVAN_DIR`

```bash
  # override the scaffold location with `CARAVAN_DIR`
  CARAVAN_DIR=/tmp/caravan/cluster ./caravan cluster up
```


- The two compute nodes advertise `gpu:4` each as **fake, count-only GPUs**
- Real GPU scheduling, no hardware needed (no `nvidia-smi` telemetry) by default.

#### 🎮 Using a real GPU

If you have an NVIDIA GPU with Docker GPU support enabled, opt in with
`CARAVAN_GPU=real`: the `c1` compute node gets the real GPU passed through
(via NVIDIA Container Toolkit, no CUDA base image needed), while `c2` keeps
the fake GPUs since most machines only have one physical GPU to give.

🐧 **On Linux / Windows (Git Bash)**

```bash
  CARAVAN_GPU=real ./caravan cluster up
  CARAVAN_GPU=real ./caravan submit workloads/gpu_example.sh
```

⊞ **On Windows (PowerShell)**

```powershell
  $env:CARAVAN_GPU = "real"
  ./caravan.exe cluster up
  ./caravan.exe submit workloads/gpu_example.sh
```

> [!NOTE]
> **macOS isn't supported here.** Docker Desktop for Mac has no GPU
> passthrough mechanism at all, and Macs don't have NVIDIA GPUs (Apple
> Silicon uses its own GPU; even older Intel MacBooks shipped AMD, not
> NVIDIA). Setting `CARAVAN_GPU=real` there fails `c1` outright with
> something like *"could not select device driver 'nvidia' with
> capabilities: [[gpu]]"*. Tested combination is an NVIDIA GPU + Docker
> Desktop on Windows/WSL2 or native Linux. The default fake-GPU cluster
> is unaffected and works the same everywhere.

---

### 📋 Checking Cluster Status

Print container state, then `sinfo`

```bash
  # check job status
  ./caravan cluster status
  ```

<details>
<summary>Output:</summary>

    hitesh@Mac Caravan % ./caravan cluster status
    CONTAINER ID  IMAGE                           COMMAND     CREATED         STATUS                 PORTS       NAMES
    ce05d9f1bbd9  localhost/caravan-slurm:latest  slurmd      15 seconds ago  Up 1 second (healthy)              c1
    0ef798e4b97c  localhost/caravan-slurm:latest  slurmd      15 seconds ago  Up 1 second (healthy)              c2
    2f3e0982f69c  localhost/caravan-slurm:latest  slurmctld   14 seconds ago  Up 14 seconds                      slurmctld

    PARTITION AVAIL  TIMELIMIT  NODES  STATE NODELIST
    gpu*         up   infinite      2    unk c[1-2]
                                c2
</details>


### Check logs slurmctld squeue

 **With Podman**

```bash
  podman exec slurmctld squeue
```

<details>
<summary>Output:</summary>

    hitesh@Mac Caravan %   podman exec slurmctld squeue
    JOBID PARTITION     NAME     USER ST       TIME  NODES NODELIST(REASON)
    1       gpu caravan-     root PD       0:00      1 (Nodes required for job are DOWN, DRAINED or reserved for jobs in higher priority partitions)

</details>


**With Docker**

```bash
  docker exec slurmctld squeue
  
  # download slurmctld logs
  docker exec slurmctld cat /var/log/slurm/slurmctld.log
  docker exec c1 cat /var/log/slurm/slurmd.log
```


<details>
<summary>Output:</summary>

    [2026-06-27T19:31:38.456] Launching batch job 1 for UID 0
    [2026-06-27T19:31:38.464] CPU frequency setting not configured for this node
    [2026-06-27T19:31:38.573] [1.batch] done with step
    [2026-06-27T19:36:33.488] Launching batch job 2 for UID 0
    [2026-06-27T19:36:33.502] CPU frequency setting not configured for this node
    [2026-06-27T19:36:33.700] [2.batch] done with step
    [2026-06-27T19:38:48.497] Launching batch job 3 for UID 0
    [2026-06-27T19:38:48.513] CPU frequency setting not configured for this node
    [2026-06-27T19:38:48.725] [3.batch] done with step
    [2026-06-27T19:39:18.499] Launching batch job 4 for UID 0
    [2026-06-27T19:39:18.508] CPU frequency setting not configured for this node
    [2026-06-27T19:39:18.609] [4.batch] done with step

</details>



---

### 📥 Submitting a HPC Workload

Create a simple job script (see `workloads/submit_example.sh`) and submit it:

```bash
  # submit the example script
  ./caravan submit workloads/submit_example.sh

```

<details>
<summary>Output</summary>

    → submitting workloads/submit_example.sh to local Slurm cluster in /Users/hitesh/.caravan/cluster
    #!/usr/bin/env bash
    #SBATCH --job-name=caravan-test
    #SBATCH --output=caravan-test.out
    #SBATCH --time=00:01:00
    #SBATCH --ntasks=1

    echo "Hello from Caravan job on $(hostname)"
    sleep 5
    echo "Done"

</details>

For a real GPU instead (see "Using a real GPU" above), submit
`workloads/gpu_example.sh` the same way under `CARAVAN_GPU=real`: it
requests `--gres=gpu:1` and prints actual `nvidia-smi` output.

---

### 💥 Tear Down Slurm Cluster
Tear down cluster and option to clean mounted volumes as well

```bash  
  ./caravan cluster down      
  ./caravan cluster down  -v  # Also remove disc volumes  
```  

---

## 👨‍💻 Development

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
  ├── embed.go                 # go:embed slurm-cluster/* (must live next to it)
  ├── internal/
  │   ├── cli/                 # cobra commands
  │   │   ├── root.go
  │   │   ├── cluster.go       # caravan cluster up|down|status
  │   │   └── submit.go        # caravan submit <script.sh>
  │   └── cluster/
  │       ├── engine.go        # Engine interface + DockerEngine/PodmanEngine
  │       ├── compose.go       # compose file paths + CARAVAN_GPU gate
  │       ├── extract.go       # Scaffold (wired from embed.go) + extraction
  │       ├── status.go        # Up/Down/Status
  │       ├── submit.go        # Submit — streams a script into sbatch
  │       └── util.go          # process-running helpers
  ├── slurm-cluster/           # the GPU Slurm cluster, embedded in the binary
  │   ├── Dockerfile · entrypoint.sh
  │   ├── docker-compose.yml · slurm.conf · gres.conf · cgroup.conf
  │   └── docker-compose.gpu.yml · slurm.gpu.conf · gres.gpu.conf  # CARAVAN_GPU=real overlay
  └── workloads/                # example job scripts
      ├── submit_example.sh
      └── gpu_example.sh
```


----

## 🗺️ Roadmap

Caravan grows from "runs a cluster" to "runs your work on it."

- **cluster** *(here)* — `up` / `down` / `status`, behind a single `Engine` interface
  implemented by Docker and Podman today; cloud / bare-metal backends later.
- **submit** *(here)* — `caravan submit script.sh` streams the script straight into
  `sbatch` on the controller. `logs` to follow it; recording/rerun next.
- **rerun** — re-launch a past job by id, reproducibly.
- **exp** — group runs into experiments and compare them.

Each new capability sits behind a `Backend` interface (Slurm today), so swapping
the execution target later doesn't touch the CLI above it.

---

## License
*© 2026 [Hitesh Kumar Sahu](https://hiteshsahu.com) · Licensed under [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0)*

