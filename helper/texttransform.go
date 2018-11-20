package helper

import (
	"strings"
)

func TextTransform(message string) string {
	s := strings.Split(message, ";")
	return s[0]
}
