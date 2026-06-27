package cli

import (
	"github.com/hiteshsahu/caravan/internal/cluster"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit <script.sh>",
	Short: "Submit a Slurm workload script to the local cluster",
	Args:  cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return cluster.Submit(args[0])
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
