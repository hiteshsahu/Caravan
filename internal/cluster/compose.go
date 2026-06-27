package cluster

import (
	"os"
	"path/filepath"
)

const project = "caravan"

func composeFile(dir string) string {
	return filepath.Join(dir, "docker-compose.yml")
}

func composeGPUFile(dir string) string {
	return filepath.Join(dir, "docker-compose.gpu.yml")
}

// GPUEnabled reports whether real GPU passthrough was requested via
// CARAVAN_GPU=real. Default is the fake, count-only GPUs that need no
// hardware.
func GPUEnabled() bool {
	return os.Getenv("CARAVAN_GPU") == "real"
}
