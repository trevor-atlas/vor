package utils

import (
	"os"
	"strings"
)

func GetEnv(name string) string {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if pair[0] == name {
			return pair[1]
		}
	}
	return ""
}

func LeftPad(s string, padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen) + s
}
