package stringutil

import "strings"

// IsAlpha checks if a given character is alphanumeric
func IsAlpha(character byte) bool {
	var result bool
	result = character > 47 && character < 58
	return result
}

// IsCharacter checks if a given character is an character (letter)
func IsCharacter(character byte) bool {
	var result bool

	result = (character > 64 && character < 91 ||
		character > 96 && character < 123)
	return result
}

// IsUppercase checks if a given character is upper case or not
func IsUppercase(character byte) bool {
	var result bool
	result = character > 64 && character < 91
	return result
}

// IsWhitespace checks if a given character is a whitespace
func IsWhitespace(character byte) bool {
	var result bool
	result = (character == ' ' ||
		character == '\n' ||
		character == '\r' ||
		character == '\v' ||
		character == '\f' ||
		character == '\t')
	return result
}

// BreakBySpaces shorthand for BreakByToken(string, ' ')
func BreakBySpaces(s string) (string, string) {
	return BreakByToken(s, ' ')
}

// BreakByToken breaks a string into two parts only if the specified token
// was found. Returns the left hand side of the string as first return value
// and the remainder of the string as the second.
// If token wansn't found then returns the input string as first
// and empty string as second.
func BreakByToken(s string, tok byte) (string, string) {
	s = strings.TrimSpace(s)
	size := len(s)
	index := 0
	for {
		if index >= size {
			break
		}

		if s[index] == tok {
			return s[0:index], EatSpaces(s[index+1:])
		}

		index++
	}

	return s, ""
}

// EatUntilAlpha eats everything until an alphanumeric is found
func EatUntilAlpha(s string) string {
	index := 0
	size := len(s)
	for {
		if index >= size {
			break
		}

		if IsAlpha(s[index]) {
			return s[index:]
		}

		index++
	}
	return s
}

// EatSpaces eats every space to left side of a string
func EatSpaces(s string) string {
	index := 0
	size := len(s)
	for {
		if index >= size {
			break
		}

		if s[index] != ' ' {
			return s[index:]
		}

		index++
	}

	return s
}

// SubstringBetween retrieves string between two specified tokens.
// For example: A call like SubstringBetween("My string *bold*", '*', '*')
// retrieves the first string found between the asterisks. In this case
// the returned string would be bold.
func SubstringBetween(s string, openTok, closeTok byte) string {
	_, remainder := BreakByToken(s, openTok)
	result, _ := BreakByToken(remainder, closeTok)
	return result
}

// Advance advances the given string by the given N amount
func Advance(s *string, n int) {
	if n == 0 {
		return
	}

	size := len(*s)
	if n < 0 || n > size {
		return
	}

	*s = (*s)[n:]
}
