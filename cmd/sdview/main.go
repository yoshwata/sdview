package main

import (
	"os"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/yoshwata/sdview/pkg/cmd"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = cmd.NewCmdLab(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
)

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	flags := pflag.NewFlagSet("kubectl-ns", pflag.ExitOnError)
	pflag.CommandLine = flags

	rootCmd := cmd.NewCmdLab(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
