package main

import "embed"

// The GPU Slurm cluster scaffold travels inside the binary, embedded here
// since go:embed directives can't reach outside their own directory tree.
//
//go:embed slurm-cluster/*
var slurmCluster embed.FS
