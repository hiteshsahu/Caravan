package main

import (
	"fmt"
	"os"

	"github.com/hiteshsahu/caravan/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "caravan:", err)
		os.Exit(1)
	}
}
