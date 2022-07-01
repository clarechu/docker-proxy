package cmd

import (
	"github.com/clarechu/docker-proxy/example"
	"github.com/clarechu/docker-proxy/pkg/models"
	"github.com/clarechu/docker-proxy/pkg/router"
	"github.com/spf13/cobra"
)

// GetServerCmd returns the root of the cobra command-tree.
func GetServerCmd(args []string) *cobra.Command {

	root := &models.Root{
		App: example.NewApp1(),
	}
	rootCmd := &cobra.Command{
		Use:   "server",
		Short: "start docker proxy server",
		Long: `
Tips  Find more information at: https://github.com/clarechu/docker-proxy
Example:
docker-proxy server
`,
		Run: func(cmd *cobra.Command, args []string) {
			server := router.NewServer(root)
			server.Run()

		},
	}
	addFlag(rootCmd, root)
	rootCmd.AddCommand(VersionCommand())
	return rootCmd
}

func addFlag(rootCmd *cobra.Command, root *models.Root) {
	rootCmd.PersistentFlags().Int32Var(&root.Port, "port", 7777, "proxy server ports")
}
