// Copyright Jetstack Ltd. See LICENSE for details.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	tarmakv1alpha1 "github.com/jetstack/tarmak/pkg/apis/tarmak/v1alpha1"
)

var globalFlags = &tarmakv1alpha1.Flags{}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "tarmak",
	Short: "Tarmak is a toolkit for provisioning and managing Kubernetes clusters.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(args []string) {
	RootCmd.SetArgs(args)

	// evalutate command that is gonna be run
	command, commandArgs, err := RootCmd.Traverse(args)

	// escape pass through commands (kubectl, ssh) if necessary
	if err == nil && (command.Use == "kubectl" || command.Use == "ssh") {
		// if no escape exists already add one
		if !stringSliceContains(commandArgs, "--") {
			pos := len(args) - len(commandArgs)
			newArgs := append(args[:pos], append([]string{"--"}, args[pos:]...)...)
			RootCmd.SetArgs(newArgs)
			// this line helps debugging fmt.Printf("rewriting args\noriginal args=%v\ncommand  args=%v\nnew      args=%v\n", args, commandArgs, newArgs)
		}
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func stringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func init() {
	// set default tarmak config folder
	tarmakConfigPath := "~/.tarmak"
	if envConfigPath := os.Getenv("TARMAK_CONFIG"); envConfigPath != "" {
		tarmakConfigPath = envConfigPath
	}

	RootCmd.PersistentFlags().StringVarP(
		&globalFlags.ConfigDirectory,
		"config-directory",
		"c",
		tarmakConfigPath,
		"config directory for tarmak's configuration",
	)

	RootCmd.PersistentFlags().BoolVarP(
		&globalFlags.Verbose,
		"verbose",
		"v",
		false,
		"enable verbose logging",
	)

	RootCmd.PersistentFlags().BoolVar(
		&globalFlags.KeepContainers,
		"keep-containers",
		false,
		"do not clean-up terraform/packer containers after running them",
	)
}
