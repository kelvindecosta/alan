package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kelvindecosta/alan/parts"
	"github.com/kelvindecosta/alan/utils"
	"github.com/spf13/cobra"
)

var designCmd = &cobra.Command{
	Use:   "design machine",
	Short: "Design definition of machine",
	RunE:  design,
	Args:  cobra.ExactArgs(1),
}

var description, states, symbols, blank, alphabet, start, final string
var noFinalState bool
var addTransitions, remTransitions []string

func init() {
	designCmd.Flags().StringVarP(&description, "description", "D", "", "description of machine")
	designCmd.Flags().StringVarP(&states, "states", "Q", "", "states of machine")
	designCmd.Flags().StringVarP(&symbols, "symbols", "L", "", "symbols of tape")
	designCmd.Flags().StringVarP(&blank, "blank", "B", "", "blank tape symbol")
	designCmd.Flags().StringVarP(&alphabet, "alphabet", "A", "", "alphabet of input")
	designCmd.Flags().StringVarP(&start, "start", "S", "", "start state of machine")
	designCmd.Flags().StringVarP(&final, "final", "F", "", "set of final states")
	designCmd.Flags().BoolVar(&noFinalState, "no-final-state", false, "whether set of final states is empty")
	designCmd.Flags().StringSliceVar(&addTransitions, "add-trans", nil, "set of transitions to be added   : <q1><s1><s2><q2><L/R> separated by ,")
	designCmd.Flags().StringSliceVar(&remTransitions, "rem-trans", nil, "set of transitions to be removed : <q1><s1><s2><q2><L/R> separated by ,")
	RootCmd.AddCommand(designCmd)
}

func design(cmd *cobra.Command, args []string) error {
	_, err := utils.CreateFileIfNotExist(args[0])
	if err != nil {
		return err
	}

	tempFile, err := ioutil.TempFile(".", "")
	if err != nil {
		return err
	}

	defer os.Remove(tempFile.Name())

	d := parts.Definition{}

	// Read file
	bytes, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	if len(bytes) > 0 {
		// Unmarshal JSON
		err = d.UnmarshalJSON(bytes)
		if err != nil {
			panic(err)
		}
	}

	if description != "" {
		d.Description = description
	}

	if states != "" {
		d.States = states
	}

	if symbols != "" {
		d.Symbols = symbols
	}

	if blank != "" {
		d.Blank = blank
	}

	if alphabet != "" {
		d.Alphabet = alphabet
	}

	if start != "" {
		d.Start = start
	}

	if final != "" && noFinalState {
		return fmt.Errorf("set of finals states is empty but also '%s'", final)
	}

	if noFinalState {
		d.Final = ""
	}

	if final != "" {
		d.Final = final
	}

	for _, t := range addTransitions {
		d.Transitions[t[:2]] = t[2:]
	}

	for _, t := range remTransitions {
		_, ok := d.Transitions[t[:2]]
		if ok {
			delete(d.Transitions, t[:2])
		}
	}

	bytes, err = json.MarshalIndent(d, "", "	")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(tempFile.Name(), bytes, 0666)
	if err != nil {
		return err
	}

	os.Rename(tempFile.Name(), args[0])

	return nil
}
