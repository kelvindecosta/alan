package utils

import (
	"fmt"
)

// Plural takes a num and quantity and returns its plural
func Plural(num int, quantity string) string {
	result := fmt.Sprintf("%d %s", num, quantity)
	if num != 1 {
		result = result + "s"
	}
	return result
}
