package utils

import "strings"

func PadOutput(padding int) func(string) string {
	return func(str string) string {
		return LeftPad(str, " ", padding)
	}
}

func TitleCase(s string) string {
	r := strings.NewReplacer(
		"-", " ",
		"/", " ",
		"\n", "",
		"\t", " ")
	var result string
	for _, word := range strings.Split(r.Replace(s), " ") {
		result += strings.ToUpper(string(word[0])) + string(word[1:])
	}
	return result
}

func KebabCase(s string) string {
	r := strings.NewReplacer(
		" ", "-",
		":", "",
		"/", "-",
		"\n", "",
		"\t", "-")
	return r.Replace(s)
}

func LowerKebabCase(s string) string {
	return strings.ToLower(KebabCase(s))
}

func Contains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

// LeftPad pad s with pLen number of padStr
func LeftPad(s string, padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen) + s
}
