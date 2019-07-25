package machine

import (
	"container/list"
	"fmt"
	"strings"
)

type void struct{}
type changes struct {
	symbol    rune
	direction bool
	state     string
}

// A Machine represents the implementation for the Turing Machine datastructure
type Machine struct {
	// Formal definition
	symbols     map[rune]void
	blank       rune
	states      map[string]void
	start       string
	end         map[string]void
	transitions map[string]map[rune]changes

	// Implementation variables
	current string
	tape    *list.List
	head    *list.Element
}

// NewMachine initializes and returns a new Machine variable
func NewMachine() *Machine {
	m := Machine{}
	m.symbols = make(map[rune]void)
	m.states = make(map[string]void)
	m.end = make(map[string]void)
	m.transitions = make(map[string]map[rune]changes)

	return &m
}

// AddSymbol adds a symbol to the Machine
func (m *Machine) AddSymbol(symbol rune) {
	m.symbols[symbol] = void{}
}

// SetBlankSymbol sets a symbol as the blank symbol of the Machine
// and returns an error if the blank symbol is already set
func (m *Machine) SetBlankSymbol(symbol rune) error {
	m.AddSymbol(symbol)
	if m.blank == 0 {
		m.blank = symbol
		return nil
	}
	return fmt.Errorf("Machine already has a blank symbol : %#v", m.blank)
}

// GetSymbols returns the symbols and blank symbol of the Machine
// and an error if the Machine has no symbols
func (m *Machine) GetSymbols() ([]rune, rune, error) {
	symbols := []rune{}
	if len(m.symbols) == 0 {
		return symbols, 0, fmt.Errorf("Machine has no symbols")
	}

	for r := range m.symbols {
		symbols = append(symbols, r)
	}
	return symbols, m.blank, nil
}

// AddState adds a state to the Machine
func (m *Machine) AddState(state string) {
	m.states[state] = void{}
}

// SetStartState sets a state as the start state of the Machine
// and returns  an error if the start state is already set
func (m *Machine) SetStartState(state string) error {
	m.AddState(state)
	if m.start == "" {
		m.start = state
		return nil
	}
	return fmt.Errorf("Machine already has a start state : %s", m.start)
}

// AddEndState adds a state to the set of end states
func (m *Machine) AddEndState(state string) {
	m.AddState(state)
	m.end[state] = void{}
}

// GetStates returns the states, start state and end states of a Machine
// and an error if the Machine has no states
func (m *Machine) GetStates() ([]string, string, []string, error) {
	states := []string{}
	end := []string{}
	if len(m.states) == 0 {
		return states, "", end, fmt.Errorf("Machine has no states")
	}

	for k := range m.states {
		states = append(states, k)
	}

	for k := range m.end {
		end = append(end, k)
	}

	return states, m.start, end, nil
}

// AddTransition adds a state transition to the Machine
func (m *Machine) AddTransition(curState string, curSym, nextSym rune, dir bool, nextState string) {
	if _, ok := m.transitions[curState]; !ok {
		m.transitions[curState] = make(map[rune]changes)
	}
	m.AddState(curState)
	m.AddSymbol(curSym)
	m.AddSymbol(nextSym)
	m.AddState(nextState)

	m.transitions[curState][curSym] = changes{nextSym, dir, nextState}
}

// String returns a human readable form of the Machine
func (m Machine) String() string {
	var builder strings.Builder
	builder.WriteString("Machine\n")

	symbols, blank, _ := m.GetSymbols()
	builder.WriteString("\tSymbols : [")
	for _, r := range symbols {
		builder.WriteString(fmt.Sprintf("'%c', ", r))
	}
	builder.WriteString("]\n")
	builder.WriteString(fmt.Sprintf("\tBlank : '%c'\n", blank))

	states, start, end, _ := m.GetStates()
	builder.WriteString(fmt.Sprintf("\tStates : %v\nStart : %s\n\tEnd : %v\n", states, start, end))
	builder.WriteString(fmt.Sprintf("\tTransitions : %v\n", m.transitions))

	return builder.String()
}

// GetTape returns the tape of the Machine and the index of its head
// and an error if the Machine is not given an input tape
func (m *Machine) GetTape() (string, int, error) {
	if m.tape.Len() == 0 {
		return "", -1, fmt.Errorf("Machine is not given an input tape")
	}
	var builder strings.Builder
	i, h := 0, 0
	for e := m.tape.Front(); e != nil; e = e.Next() {
		r := e.Value.(rune)
		builder.WriteRune(r)
		if m.head.Value.(rune) == r {
			h = i
		}
		i++
	}
	return builder.String(), h, nil
}

// Reset resets the runtime state of the Machine by loading the tape with an input string
func (m *Machine) Reset(input string) {
	m.current = m.start
	m.tape = list.New()
	for _, r := range input {
		m.tape.PushBack(r)
	}
	m.head = m.tape.Front()
}

// Step implements one iteration of the computation by the Machine
// It returns whether or not the Machine halted
func (m *Machine) Step() bool {
	next, ok := m.transitions[m.current][m.head.Value.(rune)]
	if !ok {
		return true
	}

	// Change symbol at head
	m.head.Value = next.symbol

	// Move tape
	if next.direction {
		if temp := m.head.Prev(); temp != nil {
			// Moves head to the left
			m.head = temp
		} else {
			// Moves head to the left after inserting a blank
			m.tape.PushFront(m.blank)
			m.head = m.tape.Front()
		}
	} else {
		if temp := m.head.Next(); temp != nil {
			// Moves head to the right
			m.head = temp
		} else {
			// Moves head to the right after inserting a blank
			m.tape.PushBack(m.blank)
			m.head = m.tape.Back()
		}
	}

	// Change state
	m.current = next.state

	return false
}

// Compute implements maxSteps iterations of computation by a Machine on an input string
// or until the Machine halts. It returns whether or not the Machine halted and accepted the string
func (m *Machine) Compute(input string, maxSteps uint) (bool, bool) {
	m.Reset(input)
	var steps uint
	halt := false

	for steps < maxSteps && !halt {
		steps++
		halt = m.Step()
	}

	_, accepted := m.end[m.current]
	return halt, halt && accepted
}
