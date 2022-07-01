package root

import (
	command "github.com/clarechu/docker-proxy/cmd"
	"github.com/spf13/cobra"
)

// GetRootCmd returns the root of the cobra command-tree.
func GetRootCmd(args []string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "docker-proxy",
		Short: "docker-proxy ...",
		Long: `
Tips  Find more information at: https://github.com/clarechu/docker-proxy
Example:
`,
	}
	rootCmd.AddCommand(command.GenerateCommand())
	rootCmd.AddCommand(command.GetServerCmd(args))
	rootCmd.AddCommand(command.VersionCommand())
	return rootCmd
}
