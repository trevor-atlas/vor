package utils

import (
	"strings"
	"unicode/utf8"
)

func PadOutput(padding int) func(string) string {
	return func(str string) string {
		return LeftPad(str, " ", padding)
	}
}

func SanitizeWhitespace(str string) string {
	r := strings.NewReplacer("\n", "", "\t", "")
	return r.Replace(str)
}

func Slashify(str string) string {
	r := strings.NewReplacer(" ", "/")
	return r.Replace(str)
}

func Dashify(str string) string {
	r := strings.NewReplacer(" ", "-")
	return r.Replace(str)
}

func StringPipe(fns ...func(string) string) func (string) string {
	return func (str string) string {
		for _, fun := range fns {
			str = fun(str)
		}
		return str
	}
}

// Wordwrap wraps text at the specified column lineWidth on word breaks
func Wordwrap(text string, lineWidth int) string {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return text
	}
	wrapped := words[0]
	spaceLeft := lineWidth - utf8.RuneCountInString(wrapped)
	for _, word := range words[1:] {
		if utf8.RuneCountInString(word) + 1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - utf8.RuneCountInString(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + utf8.RuneCountInString(word)
		}
	}

	return wrapped
}

func Capitalize(s string) string {
	return strings.ToUpper(string(s[0])) + string(s[1:])
}

func TitleCase(s string) string {
	r := strings.NewReplacer(
		"-", " ",
		"/", " ",
		"\n", "",
		"\t", "")
	var result string
	for _, word := range strings.Split(r.Replace(s), " ") {
		result += strings.ToUpper(string(word[0])) + string(word[1:] + " ")
	}
	return result
}

func KebabCase(s string) string {
	r := strings.NewReplacer(
		" ", "-",
		":", "",
		"/", "-",
		"\n", "",
		"\t", "-",
		"&", "and",
		",", "",
		".", "")
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

