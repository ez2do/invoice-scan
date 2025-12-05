package pkg

import (
	"strings"
	"unicode"
)

// StringTrimSpace -- trim space of string
func StringTrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// IsStringEmpty -- check if string is empty
func IsStringEmpty(s string) bool {
	return s == ""
}

// IsStringNotEmpty -- check if string is not empty
func IsStringNotEmpty(s string) bool {
	return s != ""
}

// UCFirst uppercase first character
func UCFirst(s string) string {
	var first rune
	for _, v := range s {
		first = v
		break
	}

	if unicode.IsUpper(first) {
		return s
	}
	return string(unicode.ToUpper(first)) + s[1:]
}

func caseHelper(s string, isCamel bool) []string {
	replacer := strings.NewReplacer("_", " ", "-", " ", ".", " ")
	raw := replacer.Replace(s)
	words := strings.Fields(raw)
	if !isCamel {
		var camelWords []string
		for _, word := range words {
			var prev rune
			var prevIdx = 0
			for i, c := range word {
				if i == 0 {
					prev = c
					continue
				}
				if i-prevIdx > 1 {
					if unicode.IsUpper(c) && unicode.IsLower(prev) {
						camelWords = append(camelWords, word[prevIdx:i])
						prevIdx = i
					} else if unicode.IsLower(c) && unicode.IsUpper(prev) {
						camelWords = append(camelWords, word[prevIdx:i-1])
						prevIdx = i - 1
					}
				} else {
					if i == len(word)-1 {
						if unicode.IsUpper(c) && unicode.IsLower(prev) {
							camelWords = append(camelWords, word[prevIdx:i])
							prevIdx = i
						}
					}
				}
				prev = c
			}

			if prevIdx < len(word) {
				camelWords = append(camelWords, word[prevIdx:])
			}
		}
		return camelWords
	}
	return words
}

// CamelCase converts string s to camel case
func CamelCase(s string) string {
	words := caseHelper(s, true)
	for idx, w := range words {
		if idx > 0 {
			words[idx] = UCFirst(w)
		}
	}

	return strings.Join(words, "")
}

// SnakeCase converts string s to snake case
func SnakeCase(s string) string {
	words := caseHelper(s, false)
	for idx, w := range words {
		words[idx] = strings.ToLower(w)
	}

	return strings.Join(words, "_")
}

// KebabCase converts string s to kebab case
func KebabCase(s string) string {
	words := caseHelper(s, false)
	for idx, w := range words {
		words[idx] = strings.ToLower(w)
	}

	return strings.Join(words, "-")
}
