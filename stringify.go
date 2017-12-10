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
func (lx LamExpr) String() string {
	return lx.deDeBruijn([]string{}, 0)
}

func (lx LamExpr) deDeBruijn(boundLetters []string, nextletter int) string {
	result := ""

	for _, part := range lx {
		switch part := part.(type) {
		case int:
			if part < int(len(boundLetters)) && boundLetters[part] != "" {
				result += boundLetters[part] + " "
			} else {
				newLetter := intToLetter(nextletter)
				nextletter++

				for i := int(len(boundLetters)); i < part; i++ {
					boundLetters = append(boundLetters, "")
				}

				boundLetters = append(boundLetters, newLetter)
				result += newLetter + " "
			}

		case LamAbst:
			if len(lx) == 1 {
				result += part.deDeBruijn(boundLetters, nextletter)
			} else {
				result = result + "(" + part.deDeBruijn(boundLetters, nextletter) + ") "
			}

		case LamExpr:
			result = result + "(" + part.deDeBruijn(boundLetters, nextletter) + ") "

		default:
			panic("invalid type in LamExpr")
		}
	}

	return strings.TrimSuffix(result, " ")
}

// String returns the lambda abstraction as a string
func (lf LamAbst) String() string {
	return lf.deDeBruijn([]string{}, 0)
}

func (lf LamAbst) deDeBruijn(boundLetters []string, nextletter int) string {
	// First make the first character undefined (for now)
	newLetter := intToLetter(nextletter)
	nextletter++

	boundLetters = append([]string{newLetter}, boundLetters...)
	result := "\\" + newLetter + "."

	lx := LamExpr(lf)
	result += lx.deDeBruijn(boundLetters, nextletter)

	return result
}
