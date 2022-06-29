package main

import (
	"flag"
	"github.com/clarechu/docker-proxy/example/root"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
	"os"
)

func init() {
	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

}

func main() {
	rootCmd := root.GetRoot1Cmd(os.Args[1:])
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
