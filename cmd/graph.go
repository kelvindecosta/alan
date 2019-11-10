package cmd

import (
	"fmt"

	"github.com/kelvindecosta/alan/machine"
	"github.com/spf13/cobra"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph file",
	Short: "Graph a Turing Machine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m := machine.NewMachine()
		m.Parse(args[0])
		fmt.Print(m.Graph())
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}
