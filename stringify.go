package LamCalc

import (
	"strconv"
	"strings"
)

func varToLetter(num LamVar) string {
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
	return "x" + strconv.Itoa(int(num)-25)
}

// String returns the Lambda Expression as a string
func (lx LamExpr) String() string {
	return lx.deDeBruijn([]string{}, 0)
}

func (lx LamExpr) deDeBruijn(boundLetters []string, nextletter LamVar) string {
	result := ""

	for _, part := range lx {
		switch part := part.(type) {
		case LamVar:
			result += part.deDeBruijn(boundLetters, nextletter) + " "

		case LamAbst:
			if len(lx) == 1 {
				result += part.deDeBruijn(boundLetters, nextletter)
			} else {
				result += "(" + part.deDeBruijn(boundLetters, nextletter) + ") "
			}

		case LamExpr:
			result += "(" + part.deDeBruijn(boundLetters, nextletter) + ") "
		}
	}

	return strings.TrimSuffix(result, " ")
}

// String returns the lambda abstraction as a string
func (la LamAbst) String() string {
	return la.deDeBruijn([]string{}, LamVar(0))
}

func (la LamAbst) deDeBruijn(boundLetters []string, nextletter LamVar) string {
	// First make the first character undefined (for now)
	newLetter := varToLetter(nextletter)
	nextletter++

	boundLetters = append([]string{newLetter}, boundLetters...)
	result := "λ" + newLetter + "."

	lx := LamExpr(la)
	result += lx.deDeBruijn(boundLetters, nextletter)

	return result
}

// String returns the lambda variable as a string
func (lv LamVar) String() string {
	return lv.deDeBruijn([]string{}, LamVar(0))
}

func (lv LamVar) deDeBruijn(boundLetters []string, nextletter LamVar) string {
	if lv < LamVar(len(boundLetters)) && boundLetters[lv] != "" {
		return boundLetters[lv]
	}

	newLetter := varToLetter(nextletter)
	nextletter++

	for i := LamVar(len(boundLetters)); i < lv; i++ {
		boundLetters = append(boundLetters, "")
	}

	boundLetters = append(boundLetters, newLetter)
	return newLetter
}
