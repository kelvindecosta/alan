package parts

import (
	"encoding/json"
	"io/ioutil"
	"strings"
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

// Reset resets the Machine to the initial state and loads input tape
func (m *Machine) Reset(input string) {
	// Set machine on start state
	m.State = m.Rules.Start

	// Load input and set head
	m.Tape = input
	if m.Tape == "" {
		m.Tape = m.Rules.Blank
	}
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
func (m *Machine) Compute(input string, maxSteps uint) (string, bool, bool) {
	m.Reset(input)
	var steps uint
	halt := false

	for steps < maxSteps && !halt {
		steps++
		halt = m.Step()
	}

	accepted := halt && strings.Contains(m.Rules.Final, m.State)
	return m.Tape, halt, accepted
}
