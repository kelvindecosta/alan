package machine

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

type void struct{}
type changes struct {
	symbol    rune
	direction bool
	state     string
}

// Machine defines the basic data structure for a Turing machine
type Machine struct {
	symbols     map[rune]void
	blankSymbol rune
	states      map[string]void
	startState  string
	endStates   map[string]void
	transitions map[string]map[rune]changes

	currentState string
	tape         *list.List
	head         *list.Element
}

// NewMachine returns a Machine instance
func NewMachine() *Machine {
	m := Machine{}
	m.symbols = make(map[rune]void)
	m.states = make(map[string]void)
	m.endStates = make(map[string]void)
	m.transitions = make(map[string]map[rune]changes)

	return &m
}

// SetSymbol stores a symbol in the Machine definition
// The symbol can be set as the blank symbol only if it is being set for the first time
// The program halts if a symbol is set as the blank symbol when the blank symbol is already set
func (m *Machine) SetSymbol(symbol rune, isBlank bool) {
	m.symbols[symbol] = void{}
	if isBlank {
		if m.blankSymbol == 0 {
			m.blankSymbol = symbol
		} else {
			log.Fatalf("Machine got blank symbol '%c' which is already set to '%c'", symbol, m.blankSymbol)
		}
	}
}

// SetState stores a state in the Machine definition
// The state can be set as the start state and/or one of the end states if it is being set for the first time
// The program halts if a state is set as the start state when the start state is already set
func (m *Machine) SetState(state string, isStart, isEnd bool) {
	m.states[state] = void{}
	if isStart {
		if m.startState == "" {
			m.startState = state
		} else {
			log.Fatalf("Machine got start state '%s' which is already set to '%s'", state, m.startState)
		}
	}

	if isEnd {
		m.endStates[state] = void{}
	}
}

// SetTransition stores a transition in the Machine definition
func (m *Machine) SetTransition(currentState string, currentSymbol, nextSymbol rune, direction bool, nextState string) {
	if _, ok := m.transitions[currentState]; !ok {
		m.transitions[currentState] = make(map[rune]changes)
	}

	m.SetState(currentState, false, false)
	m.SetState(nextState, false, false)
	m.SetSymbol(currentSymbol, false)
	m.SetSymbol(nextSymbol, false)

	m.transitions[currentState][currentSymbol] = changes{nextSymbol, direction, nextState}
}

func (m *Machine) Reset(input string) {
	m.currentState = m.startState
	m.tape = list.New()
	if input == "" {
		m.tape.PushBack(m.blankSymbol)
	} else {
		for _, r := range input {
			m.tape.PushBack(r)
		}
	}
	m.head = m.tape.Front()
}

// Step defines one transition by the Turing machine
// and returns whether or not the machine halts due to no valid transition
func (m *Machine) Step() bool {
	change, ok := m.transitions[m.currentState][m.head.Value.(rune)]
	if !ok {
		return true
	}

	m.head.Value = change.symbol

	if change.direction {
		if temp := m.head.Next(); temp != nil {
			m.head = temp
		} else {
			m.tape.PushBack(m.blankSymbol)
			m.head = m.tape.Back()
		}
	} else {
		if temp := m.head.Prev(); temp != nil {
			m.head = temp
		} else {
			m.tape.PushFront(m.blankSymbol)
			m.head = m.tape.Front()
		}
	}

	m.currentState = change.state
	return false
}

// Compute defines a computation by the Turing machine on an input string
// and returns whether the machine halted and whether it accepted the input string
func (m *Machine) Compute(input string, maxSteps uint) (bool, bool) {
	m.Reset(input)
	var steps uint
	halt := false

	for steps < maxSteps && !halt {
		steps++
		halt = m.Step()
	}

	_, accepted := m.endStates[m.currentState]
	return halt, halt && accepted
}

// Parse parses the definition in a specified filename
// and returns the corresponding Machine
func (m *Machine) Parse(filename string) {
	commentRE := regexp.MustCompile(`(#.*)`)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(data), "\n")
	for i := range lines {
		lines[i] = commentRE.ReplaceAllString(lines[i], "")
		lines[i] = strings.TrimSpace(lines[i])
	}

	blankSymbolRE := regexp.MustCompile(`^'(.)'$`)
	stateRE := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)([.*])?$`)
	transitionRE := regexp.MustCompile(`^'(.)'\s+'(.)'\s+([<>])\s+([a-zA-Z_][a-zA-Z0-9_]*)$`)
	currentState := ""

	for index, l := range lines {
		if l == "" {
			continue
		}

		bMatch := blankSymbolRE.FindStringSubmatch(l)
		if len(bMatch) > 0 {
			m.SetSymbol(rune(bMatch[1][0]), true)
			continue
		}

		sMatch := stateRE.FindStringSubmatch(l)
		if len(sMatch) > 0 {
			currentState = sMatch[1]
			isStart := false
			isEnd := false

			for _, r := range sMatch[2] {
				switch r {
				case '.':
					isEnd = true
				case '*':
					isStart = true
				}
			}

			m.SetState(currentState, isStart, isEnd)
			continue
		}

		tMatch := transitionRE.FindStringSubmatch(l)
		if len(tMatch) > 0 {
			if currentState == "" {
				break
			}

			direction := false
			switch rune(tMatch[3][0]) {
			case '<':
				direction = false
			case '>':
				direction = true
			}

			m.SetTransition(currentState, rune(tMatch[1][0]), rune(tMatch[2][0]), direction, tMatch[4])
			continue
		}

		log.Fatalf("Error at line %d\n", index+1)
	}
}

// Graph returns a definition of the Machine in Graphviz (https://www.graphviz.org/)
func (m *Machine) Graph() string {
	var builder strings.Builder
	builder.WriteString("digraph machine {\n")
	builder.WriteString("\trankdir=LR;\n")
	builder.WriteString("\tsize=\"8,5\";\n")

	builder.WriteString("\n")
	builder.WriteString("\tnode [shape = point]; 0;\n")

	for state := range m.endStates {
		builder.WriteString(fmt.Sprintf("\tnode [shape = doublecircle]; %s;\n", state))
	}

	builder.WriteString("\tnode [shape = circle];\n")
	builder.WriteString(fmt.Sprintf("\t0 -> %s;\n", m.startState))

	for currentState := range m.states {
		for currentSymbol, change := range m.transitions[currentState] {
			dir := 'L'
			if change.direction {
				dir = 'R'
			}

			builder.WriteString(fmt.Sprintf("\t%s -> %s [ label = \"'%c', '%c', '%c'\" ];\n", currentState, change.state, currentSymbol, change.symbol, dir))
		}
	}

	builder.WriteString("}\n")

	return builder.String()
}
