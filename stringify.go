package LamCalc

import (
	"bytes"
	"strconv"
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
	buffer := bytes.NewBufferString("")
	lx.deDeBruijn(buffer, new([]string), new(int))

	return buffer.String()
}

func (lx Appl) deDeBruijn(buffer *bytes.Buffer, boundLetters *[]string, nextletter *int) {
	for i, part := range lx {
		switch part := part.(type) {
		case Var:
			part.deDeBruijn(buffer, boundLetters, nextletter)

		case Abst:
			buffer.WriteRune('(')
			part.deDeBruijn(buffer, boundLetters, nextletter)
			buffer.WriteRune(')')

		case Appl:
			if i == 0 {
				part.deDeBruijn(buffer, boundLetters, nextletter)

			} else {
				buffer.WriteRune('(')
				part.deDeBruijn(buffer, boundLetters, nextletter)
				buffer.WriteRune(')')
			}
		}

		buffer.WriteRune(' ')
	}

	buffer.Truncate(buffer.Len() - 1) // Remove final space
}

// String returns the lambda abstraction as a string
func (la Abst) String() string {
	buffer := bytes.NewBufferString("")
	la.deDeBruijn(buffer, new([]string), new(int))

	return buffer.String()
}

func (la Abst) deDeBruijn(buffer *bytes.Buffer, boundLetters *[]string, nextletter *int) {
	// Remember at which character we were
	oldNextLetter := *nextletter

	// Add the localy bound letter to the list
	newLetter := intToLetter(*nextletter)
	*nextletter++
	*boundLetters = append([]string{newLetter}, *boundLetters...)

	buffer.WriteString("λ" + newLetter + ".")
	la[0].deDeBruijn(buffer, boundLetters, nextletter)

	// Remove our local naming
	*boundLetters = (*boundLetters)[1:]
	*nextletter = oldNextLetter
}

// String returns the lambda variable as a string
func (lv Var) String() string {
	buffer := bytes.NewBufferString("")
	lv.deDeBruijn(buffer, new([]string), new(int))

	return buffer.String()
}

func (lv Var) deDeBruijn(buffer *bytes.Buffer, boundLetters *[]string, nextletter *int) {
	if int(lv) < len(*boundLetters) && (*boundLetters)[lv] != "" {
		buffer.WriteString((*boundLetters)[lv])
		return
	}

	newLetter := intToLetter(*nextletter)
	*nextletter++

	for i := len(*boundLetters); i < int(lv); i++ {
		*boundLetters = append(*boundLetters, "")
	}

	*boundLetters = append(*boundLetters, newLetter)
	buffer.WriteString(newLetter)
}
