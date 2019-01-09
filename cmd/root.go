package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd is the root command
var RootCmd = &cobra.Command{
	Use:   "alan",
	Short: "A Turing Machine Simulator",
}
