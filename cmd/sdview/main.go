package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/yoshwata/sdview/pkg/cmd"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string
	version     string

	rootCmd = cmd.NewCmdSdView(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
)

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	flags := pflag.NewFlagSet("kubectl-sdview", pflag.ExitOnError)
	pflag.CommandLine = flags

	rootCmd := cmd.NewCmdSdView(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of sdview",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("version: %s\n", version)
		},
	}
	rootCmd.AddCommand((versionCmd))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
