package util

import "strings"

func ParseStyle(str string) []string {
	return strings.Split(str, ",")
}
