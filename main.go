package main

import (
	"fmt"
	"os"

	"github.com/hiteshsahu/caravan/internal/cli"
	"github.com/hiteshsahu/caravan/internal/cluster"
)

func main() {
	cluster.Scaffold = slurmCluster
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "caravan:", err)
		os.Exit(1)
	}
}
