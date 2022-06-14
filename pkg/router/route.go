package router

import (
	"github.com/ClareChu/docker-proxy/pkg/models"
	"github.com/spf13/cobra"
)

type Root struct {
	Port    int32      `json:"port"`
	Timeout int32      `json:"timeout"`
	App     models.App `json:"app"`
}

// GetRootCmd returns the root of the cobra command-tree.
func GetRootCmd(args []string) *cobra.Command {
	root := &Root{}
	rootCmd := &cobra.Command{
		Use:   "docker-proxy",
		Short: "docker-proxy ...",
		Long: `
Tips  Find more information at: https://github.com/clarechu/docker-proxy
Example:

`,
		Run: func(cmd *cobra.Command, args []string) {
			server := NewServer(root)
			server.Run()
		},
	}
	addFlag(rootCmd, root)
	rootCmd.AddCommand(VersionCommand())
	return rootCmd
}

func addFlag(rootCmd *cobra.Command, root *Root) {
	rootCmd.PersistentFlags().Int32Var(&root.Port, "port", 7777, "proxy server ports")
	rootCmd.PersistentFlags().Int32Var(&root.Timeout, "timeout", 5, "proxy server timeout")
	rootCmd.PersistentFlags().Int32Var(&root.Timeout, "timeout", 5, "proxy server timeout")

}
