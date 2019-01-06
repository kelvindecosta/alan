package parts

// Transition represents a single entry in the transition table for a Turing Machine
type Transition struct {
	CurrentState  string `json:"current-state"`
	CurrentSymbol string `json:"current-symbol"`
	NextSymbol    string `json:"next-symbol"`
	NextState     string `json:"next-state"`
	Movement      string `json:"movement"`
}
