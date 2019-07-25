package main

import (
	"fmt"

	"github.com/kelvindecosta/alan/machine"
)

func main() {
	m := machine.NewMachine()
	m.SetBlankSymbol(' ')
	m.SetStartState("A")
	m.AddEndState("G")
	m.AddTransition("A", 'X', 'X', true, "A")
	m.AddTransition("A", 'Y', 'Y', true, "A")
	m.AddTransition("A", '0', 'X', false, "B")
	m.AddTransition("A", '1', 'Y', false, "F")
	m.AddTransition("A", ' ', ' ', false, "G")
	m.AddTransition("B", '0', '0', false, "B")
	m.AddTransition("B", '1', '1', false, "B")
	m.AddTransition("B", ' ', ' ', true, "C")
	m.AddTransition("B", 'X', 'X', true, "C")
	m.AddTransition("B", 'Y', 'Y', true, "C")
	m.AddTransition("F", '0', '0', false, "F")
	m.AddTransition("F", '1', '1', false, "F")
	m.AddTransition("F", ' ', ' ', true, "E")
	m.AddTransition("F", 'X', 'X', true, "E")
	m.AddTransition("F", 'Y', 'Y', true, "E")
	m.AddTransition("C", '0', 'X', true, "D")
	m.AddTransition("C", 'X', 'X', true, "D")
	m.AddTransition("E", '1', 'Y', true, "D")
	m.AddTransition("E", 'Y', 'Y', true, "D")
	m.AddTransition("D", '0', '0', true, "D")
	m.AddTransition("D", '1', '1', true, "D")
	m.AddTransition("D", ' ', ' ', false, "A")
	m.AddTransition("D", 'X', 'X', false, "A")
	m.AddTransition("D", 'Y', 'Y', false, "A")
	m.AddTransition("G", 'X', '0', false, "G")
	m.AddTransition("G", 'Y', '1', false, "G")

	h, a := m.Compute("101", 200)
	fmt.Printf("%v, %v\n", h, a)
}
