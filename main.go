package main

import (
	_ "github.com/ClareChu/docker-proxy/pkg/proxy"
	"github.com/ClareChu/docker-proxy/pkg/router"
	"k8s.io/klog/v2"
	"os"
)

func init() {
	klog.InitFlags(nil)
	//pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

func main() {
	rootCmd := router.GetRootCmd(os.Args[1:])
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
