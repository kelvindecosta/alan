package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kelvindecosta/alan/parts"
	"github.com/kelvindecosta/alan/utils"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

var runCmd = &cobra.Command{
	Use:   "run machine",
	Short: "Run machine on input(s)",
	RunE:  run,
	Args:  cobra.ExactArgs(1),
}

var inputDirect string
var inputFilename string
var verbose bool
var maxSteps uint

func init() {
	runCmd.Flags().StringVarP(&inputDirect, "input", "i", "", "direct input")
	runCmd.Flags().StringVarP(&inputFilename, "file", "f", "", "path to input file")
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	runCmd.Flags().UintVarP(&maxSteps, "max-steps", "s", 200, "maximum steps until halt")
	runCmd.MarkFlagFilename("file")
	RootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) error {
	m := parts.Machine{}
	m.Load(args[0])

	var inputs []string

	if inputDirect == "" && inputFilename == "" {
		return errors.New("requires atleast one input source")
	}

	if inputDirect != "" {
		inputs = append(inputs, inputDirect)
	}

	if inputFilename != "" {
		file, err := os.Open(inputFilename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		s := bufio.NewScanner(file)
		for s.Scan() {
			inputs = append(inputs, s.Text())
		}
	}

	for index, tape := range inputs {
		if verbose {
			w, _ := terminal.Width()
			trace(m, tape, maxSteps)
			if index != len(inputs)-1 {
				fmt.Println("\n" + strings.Repeat("―", int(w)))
			}
		} else {
			output, halt, accepted := m.Compute(tape, maxSteps)
			var color string
			if halt {
				if accepted {
					color = "green"
				} else {
					color = "red"
				}
			} else {
				color = "blue"
			}
			fmt.Printf("%s -> %s\n", tape, ansi.Color(output, color))
		}
	}

	return nil
}

func trace(m parts.Machine, input string, maxSteps uint) {
	m.Reset(input)
	var steps uint
	halt := false

	fmt.Println(ansi.Color("\n  STEP  STATE  TAPE\n", "yellow"))
	fmt.Println(fmt.Sprintf("%5d\t%5s", steps, m.State) + "\t" + utils.Highlight(m.Tape, m.Head))

	for steps < maxSteps && !halt {
		steps++
		halt = m.Step()
		fmt.Println(fmt.Sprintf("%5d\t%5s", steps, m.State) + "\t" + utils.Highlight(m.Tape, m.Head))
	}
	var output, color string
	if halt {
		if strings.Contains(m.Rules.Final, m.State) {
			color = "green"
			output = fmt.Sprintf("\nMachine halted in %s and accepted input", utils.Plural(int(steps), "step"))
		} else {
			color = "red"
			output = fmt.Sprintf("\nMachine halted in %s and rejected input", utils.Plural(int(steps), "step"))
		}
	} else {
		color = "blue"
		output = fmt.Sprintf("\nMachine cannot decide in %s", utils.Plural(int(maxSteps), "step"))
	}

	fmt.Println(ansi.Color(output, color))
}
