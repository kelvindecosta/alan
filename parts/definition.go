package parts

import (
	"encoding/json"
)

// Definition represents the formal definition of a Turing Machine
type Definition struct {
	Description string            `json:"description"`
	States      string            `json:"states"`
	Symbols     string            `json:"symbols"`
	Blank       string            `json:"blank"`
	Alphabet    string            `json:"alphabet"`
	Start       string            `json:"start-state"`
	Final       string            `json:"final-states"`
	Transitions map[string]string `json:"-"`
}

// MarshalJSON marshals a Definition instance to JSON
func (d Definition) MarshalJSON() ([]byte, error) {
	type Alias Definition
	Transitions := []Transition{}
	for k, v := range d.Transitions {
		Transitions = append(Transitions, Transition{
			CurrentState:  string(k[0]),
			CurrentSymbol: string(k[1]),
			NextSymbol:    string(v[0]),
			NextState:     string(v[1]),
			Movement:      string(v[2]),
		})
	}

	return json.Marshal(&struct {
		*Alias
		Transitions []Transition `json:"transitions"`
	}{
		Transitions: Transitions,
		Alias:       (*Alias)(&d),
	})
}

// UnmarshalJSON unmarshals JSON to a Definition instance
func (d *Definition) UnmarshalJSON(data []byte) error {
	type Alias Definition
	aux := &struct {
		*Alias
		Transitions []Transition `json:"transitions"`
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	d.Transitions = make(map[string]string)

	for _, t := range aux.Transitions {
		key := t.CurrentState + t.CurrentSymbol
		val := t.NextSymbol + t.NextState + t.Movement
		d.Transitions[key] = val
	}
	return nil
}
