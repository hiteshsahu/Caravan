# Caravan 🐪🐫🛕

> **A CLI for GPU Slurm stand up a cluster and run workloads on it.**

![](./img/cover.jpg)

### Purpose
Slurm is powerful but a chore to operate and submit to. 

Caravan is a single binary that carries a GPU Slurm cluster *inside it* and makes it a one-liner to bring up,
submit jobs to, and (next) track experiments and rerun workloads.

Its completing project to:
- [squint](https://github.com/hiteshsahu/squint) watching the queue and 
- [gpu-lens](https://github.com/hiteshsahu/gpu-lens) watching the GPUs.

> Caravan **uses** Slurm — it doesn't replace it. Slurm stays the scheduler;

> Caravan is the control plane and developer experience around it.

---

# How to Use

## ⚡ Prerequisites

📟 Requires **Go 1.22+**

For windows
> choco install golang
More detail on this [Medium Post](https://medium.com/@hiteshkrsahu/installing-go-on-windows-the-5-minute-guide-and-the-gotchas-nobody-mentions-878eb3ea2277)

**Slurm** will be installed as Container Image

### 🐳 Container engine
Caravan works with both `Docker` or `Podman`, it auto-detects (`Docker` first, then
  `Podman`) and uses the matching Compose.

  On Podman it uses `podman-compose` if installed, otherwise it will try to use `podman compose`.


####  Using Podman

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


##  1. Get the CLI  📦

There are 2 ways to get caravan cli

### ✨ 1. Using Released Binary
Fetch release build:

```bash
  go install github.com/hiteshsahu/caravan@latest
```

### ⚙️ 2. Build Locally (for devs)

First, Build using local source:

####  On macOS

> go build -o caravan .

A `caravan` binary will be created at project root for starting the cluster.


####  On Windows
> go build -o caravan.exe .

For Windows, `caravan.exe` will be created along with `caravan` binaray


## 2. Verify

Varify if CLI is ready

### With Local Build

#### On MacOS/Linux
> ./caravan

####  On Windows 
> ./caravan.exe --help

### With Released build

> caravan


---

## 2. Start Your Caravan 🐪🐫

**Once you have CLI ready, you can start cluster and submit jobs**

![Caravan Batch](./img/illustration.jpeg)


###  Quick start with released CLI binary

**On macOS**

```bash

caravan cluster up        # build + start (controller + 2 fake-GPU nodes)
caravan cluster down      # stop (-v to also wipe volumes)
caravan cluster status    # container state + sinfo
caravan submit <script>   # stream a script into sbatch on the controller
```

### Quick start with local binary.

> Replace `caravan` with `./caravan` and you can use them for local CLI

**On macOS**

```bash

./caravan cluster up        # build + start (controller + 2 fake-GPU nodes)
./caravan cluster down      # stop (-v to also wipe volumes)
./caravan cluster status    # container state + sinfo
./caravan submit <script>   # stream a script into sbatch on the controller
```


**On Windows**

**Note:** for PowerShell users: use `./caravan.exe` instead of `./caravan`.

```bash

./caravan.exe cluster up        # build + start (controller + 2 fake-GPU nodes)
./caravan.exe cluster status    # container state + sinfo
./caravan.exe submit examples/submit_example.sh
./caravan.exe cluster down      # stop (-v to also wipe volumes)
```


---

## How it works

### ▶️ Start Cluster 

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
- The real GPU scheduling, no hardware needed (no `nvidia-smi` telemetry).

---

### 📋 Check Status
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


```bash
  podman exec slurmctld squeue
```

<details>
<summary>Output:</summary>
hitesh@Mac Caravan %   podman exec slurmctld squeue

     JOBID PARTITION     NAME     USER ST       TIME  NODES NODELIST(REASON)
         1       gpu caravan-     root PD       0:00      1 (Nodes required for job are DOWN, DRAINED or reserved for jobs in higher priority partitions)
</details>


---

### 📥 Submit a job

Create a simple job script (see `examples/submit_example.sh`) and submit it:

```bash
  # submit the example script
  ./caravan submit examples/submit_example.sh

```

<details>
<summary>Output:</summary>
  → submitting examples/submit_example.sh to local Slurm cluster in /Users/hitesh/.caravan/cluster
  #!/usr/bin/env bash
  #SBATCH --job-name=caravan-test
  #SBATCH --output=caravan-test.out
  #SBATCH --time=00:01:00
  #SBATCH --ntasks=1

echo "Hello from Caravan job on $(hostname)"
sleep 5

echo "Done"
</details>

---

### 💥 Clean up
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
  ├── internal/
  │   ├── cli/                 # cobra commands
  │   │   ├── root.go
  │   │   ├── cluster.go       # caravan cluster up|down|status
  │   │   └── submit.go        # caravan submit <script.sh>
  │   └── cluster/
  │       ├── engine.go        # Engine interface + DockerEngine/PodmanEngine
  │       ├── compose.go       # compose file path + project name
  │       ├── extract.go       # embedded scaffold + extraction
  │       ├── status.go        # Up/Down/Status
  │       ├── submit.go        # Submit — streams a script into sbatch
  │       ├── util.go          # process-running helpers
  │       └── assets/          # the GPU Slurm cluster, embedded in the binary
  │           ├── Dockerfile · entrypoint.sh
  │           ├── docker-compose.yml
  │           └── slurm.conf · gres.conf
```

The cluster definition is embedded with `//go:embed`, so the binary is
self-contained — there's no separate cluster repo to clone. `caravan cluster up`
extracts it and runs it via a `docker`/`podman` `Engine`.

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

