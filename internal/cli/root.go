package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "caravan",
	Short: "Caravan — stand up GPU Slurm clusters and run workloads on them",
	Long: "Caravan is a CLI for GPU Slurm. Today it scaffolds and runs a local Slurm\n" +
		"cluster from a scaffold embedded in the binary; submit/status/logs come next.",
	SilenceUsage: true,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
