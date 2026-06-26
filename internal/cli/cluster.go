package cli

import (
	"github.com/hiteshsahu/caravan/internal/cluster"
	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Stand up and manage the local GPU Slurm cluster",
}

var clusterUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Build and start the cluster (controller + 2 fake-GPU nodes)",
	RunE:  func(c *cobra.Command, _ []string) error { return cluster.Up() },
}

var clusterDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop the cluster",
	RunE: func(c *cobra.Command, _ []string) error {
		volumes, _ := c.Flags().GetBool("volumes")
		return cluster.Down(volumes)
	},
}

var clusterStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show container and Slurm node state",
	RunE:  func(c *cobra.Command, _ []string) error { return cluster.Status() },
}

func init() {
	clusterDownCmd.Flags().BoolP("volumes", "v", false, "also remove volumes (wipe cluster state)")
	clusterCmd.AddCommand(clusterUpCmd, clusterDownCmd, clusterStatusCmd)
	rootCmd.AddCommand(clusterCmd)
}
