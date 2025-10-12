package utils

import (
	"unicode"
)

func IsUppercase(s string) bool {
	for _, c := range s {
		if !unicode.IsUpper(c) {
			return false
		}
	}

	return true
}
