package cluster

import (
	"fmt"
	"os"
)

// Up extracts the scaffold and brings the cluster online.
func Up() error {
	dir, err := scaffoldDir()
	if err != nil {
		return err
	}
	engine, err := NewEngine()
	if err != nil {
		return err
	}
	fmt.Printf("→ engine: %s\n", engine.Name())
	if GPUEnabled() {
		fmt.Println("→ GPU: real device passthrough on c1 (CARAVAN_GPU=real)")
	}
	fmt.Printf("→ scaffolding GPU Slurm cluster in %s\n", dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if err := Extract(dir); err != nil {
		return fmt.Errorf("extract scaffold: %w", err)
	}
	fmt.Println("→ starting cluster (first run builds the image)…")
	return engine.Up(dir)
}

// Down stops the cluster; volumes also wipes its state.
func Down(volumes bool) error {
	dir, err := scaffoldDir()
	if err != nil {
		return err
	}
	engine, err := NewEngine()
	if err != nil {
		return err
	}
	return engine.Down(dir, volumes)
}

// Status prints container state, then Slurm node state.
func Status() error {
	dir, err := scaffoldDir()
	if err != nil {
		return err
	}
	engine, err := NewEngine()
	if err != nil {
		return err
	}
	_ = engine.Ps(dir)
	fmt.Println()
	return engine.Exec("slurmctld", nil, "sinfo")
}
