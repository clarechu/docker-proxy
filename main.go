package main

import (
	"github.com/clarechu/docker-proxy/example/root"
	"k8s.io/klog/v2"
	"os"
)

func init() {
	klog.InitFlags(nil)
}

func main() {
	rootCmd := root.GetRootCmd(os.Args[1:])
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
