package utils

import (
	"github.com/mgutz/ansi"
)

// Highlight highlights the text at a particular index
func Highlight(text string, index int) string {
	output := ansi.Color(string(text[index]), "cyan+i")
	if text[:index] != "" {
		output = text[:index] + output
	}
	if text[index+1:] != "" {
		output = output + text[index+1:]
	}
	return output
}
