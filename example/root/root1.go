package root

import (
	"github.com/clarechu/docker-proxy/example"
	"github.com/clarechu/docker-proxy/pkg/router"
	"github.com/spf13/cobra"
)

// GetRoot1Cmd returns the root of the cobra command-tree.
func GetRoot1Cmd(args []string) *cobra.Command {
	app := example.NewNexusApp()
	rootCmd := &cobra.Command{
		Use:   "docker-proxy",
		Short: "docker-proxy ...",
		Long: `
Tips  Find more information at: https://github.com/clarechu/docker-proxy
Example:
`,
		Run: func(cmd *cobra.Command, args []string) {
			server := router.NewNexusServer(app)
			server.Run()

		},
	}
	rootCmd.AddCommand(router.VersionCommand())
	return rootCmd
}
