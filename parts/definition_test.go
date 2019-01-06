package parts_test

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/kelvindecosta/alan/parts"
)

func TestJSONUnmarshal(t *testing.T) {
	expectedDefinition := parts.Definition{
		States:      "ABCD",
		Symbols:     "01 ",
		Blank:       " ",
		Alphabet:    "01",
		Start:       "A",
		Final:       "",
		Transitions: map[string]string{"B0": "1CL", "B1": "0BL", "B ": "1CL", "C1": "1CR", "C ": " DL", "A0": "0AR", "A1": "1AR", "A ": " BR", "C0": "0CR"},
	}

	var actualDefinition parts.Definition

	b, err := ioutil.ReadFile("testdata/definition.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, &actualDefinition)
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(actualDefinition, expectedDefinition) {
		t.Errorf("Unmarshalling Error\nExpected : %#v\nActual   : %#v", expectedDefinition, actualDefinition)
	}
}
