package LamCalc

import (
	"strconv"
	"strings"
)

func intToLetter(num int) string {
	if num < 3 {
		// x, y, z
		return string(rune(120 + num - 0))

	} else if num < 6 {
		// u, v, w
		return string(rune(117 + num - 3))

	} else if num < 26 {
		// a, b, c...
		return string(rune(97 + num - 6))
	}

	// x1, x2, x3...
	return "x" + strconv.Itoa(num-25)
}

// String returns the Lambda Expression as a string
func (lx Appl) String() string {
	return lx.deDeBruijn(new([]string), new(int))
}

func (lx Appl) deDeBruijn(boundLetters *[]string, nextletter *int) string {
	result := ""

	for i, part := range lx {
		switch part := part.(type) {
		case Var:
			result += part.deDeBruijn(boundLetters, nextletter) + " "

		case Abst:
			result += "(" + part.deDeBruijn(boundLetters, nextletter) + ") "

		case Appl:
			if i == 0 {
				result += part.deDeBruijn(boundLetters, nextletter) + " "

			} else {
				result += "(" + part.deDeBruijn(boundLetters, nextletter) + ") "
			}
		}
	}

	return strings.TrimSuffix(result, " ")
}

// String returns the lambda abstraction as a string
func (la Abst) String() string {
	return la.deDeBruijn(new([]string), new(int))
}

func (la Abst) deDeBruijn(boundLetters *[]string, nextletter *int) string {
	// Remember at which character we were
	oldNextLetter := *nextletter

	// Add the localy bound letter to the list
	newLetter := intToLetter(*nextletter)
	*nextletter++
	*boundLetters = append([]string{newLetter}, *boundLetters...)

	result := "λ" + newLetter + "."
	result += la[0].deDeBruijn(boundLetters, nextletter)

	// Remove our local naming
	*boundLetters = (*boundLetters)[1:]
	*nextletter = oldNextLetter
	return result
}

// String returns the lambda variable as a string
func (lv Var) String() string {
	return lv.deDeBruijn(new([]string), new(int))
}

func (lv Var) deDeBruijn(boundLetters *[]string, nextletter *int) string {
	if int(lv) < len(*boundLetters) && (*boundLetters)[lv] != "" {
		return (*boundLetters)[lv]
	}

	newLetter := intToLetter(*nextletter)
	*nextletter++

	for i := len(*boundLetters); i < int(lv); i++ {
		*boundLetters = append(*boundLetters, "")
	}

	*boundLetters = append(*boundLetters, newLetter)
	return newLetter
}
