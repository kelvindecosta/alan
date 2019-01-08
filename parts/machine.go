package parts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/kelvindecosta/alan/utils"
	"github.com/mgutz/ansi"
)

// Machine represents a Turing Machine
type Machine struct {
	State string
	Tape  string
	Head  int
	Rules Definition
}

// Load sets the definition of the Machine from a file
func (m *Machine) Load(definitionFile string) {
	m.Rules = Definition{}
	// Load definition
	b, err := ioutil.ReadFile(definitionFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, &(m.Rules))
	if err != nil {
		panic(err)
	}
}

// Reset resets the Machine to the initial state
func (m *Machine) Reset(input string) {
	// Set machine on start state
	m.State = m.Rules.Start

	// Load input and set head
	m.Tape = input
	m.Head = 0
}

// Step processes a single input from the tape and returns true for halt
func (m *Machine) Step() bool {
	symbol := string(m.Tape[m.Head])
	t, ok := m.Rules.Transitions[m.State+symbol]
	// Halt if no transition available
	if !ok {
		return true
	}

	// Replace symbol at head
	m.Tape = m.Tape[:m.Head] + string(t[0]) + m.Tape[m.Head+1:]

	// Change state
	m.State = string(t[1])

	// Move
	if strings.Compare(string(t[2]), "L") == 0 {
		m.Head--
	} else {
		m.Head++
	}

	if m.Head >= len(m.Tape) {
		m.Tape = m.Tape + m.Rules.Blank
	} else if m.Head == -1 {
		m.Tape = m.Rules.Blank + m.Tape
		m.Head = 0
	}

	return false
}

// Compute returns the Machine tape after processing an input, if it halted and if the input was accepted
func (m *Machine) Compute(input string, maxSteps int) (string, bool, bool) {
	m.Reset(input)
	steps := 0
	halt := false

	for steps < maxSteps && !halt {
		steps++
		halt = m.Step()
	}

	accepted := halt && strings.Contains(m.Rules.Final, m.State)
	return m.Tape, halt, accepted
}

// Trace shows step wise processing of input string by Machine
func (m *Machine) Trace(input string, maxSteps int) {
	m.Reset(input)
	steps := 0
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
			output = fmt.Sprintf("\nMachine halted in %s and accepted input", utils.Plural(steps, "step"))
		} else {
			color = "red"
			output = fmt.Sprintf("\nMachine halted in %s and rejected input", utils.Plural(steps, "step"))
		}
	} else {
		color = "blue"
		output = fmt.Sprintf("\nMachine cannot decide in %s", utils.Plural(maxSteps, "step"))
	}

	fmt.Println(ansi.Color(output, color))
}
