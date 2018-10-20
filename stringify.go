package lamcalc

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
	builder := strings.Builder{}
	lx.deDeBruijn(&builder, new([]string), new(int))

	return builder.String()
}

func (lx Appl) deDeBruijn(builder *strings.Builder, boundLetters *[]string, nextletter *int) {
	for i, part := range lx {
		switch part := part.(type) {
		case Var:
			part.deDeBruijn(builder, boundLetters, nextletter)

		case Abst:
			builder.WriteByte('(')
			part.deDeBruijn(builder, boundLetters, nextletter)
			builder.WriteByte(')')

		case Appl:
			if i == 0 {
				part.deDeBruijn(builder, boundLetters, nextletter)

			} else {
				builder.WriteByte('(')
				part.deDeBruijn(builder, boundLetters, nextletter)
				builder.WriteByte(')')
			}
		}

		// Put space between first and second part
		if i == 0 {
			builder.WriteByte(' ')
		}
	}
}

// String returns the lambda abstraction as a string
func (la Abst) String() string {
	builder := strings.Builder{}
	la.deDeBruijn(&builder, new([]string), new(int))

	return builder.String()
}

func (la Abst) deDeBruijn(builder *strings.Builder, boundLetters *[]string, nextletter *int) {
	// Remember at which character we were
	oldNextLetter := *nextletter

	// Add the localy bound letter to the list
	newLetter := intToLetter(*nextletter)
	*nextletter++
	*boundLetters = append([]string{newLetter}, *boundLetters...)

	builder.WriteRune('λ')
	builder.WriteString(newLetter)
	builder.WriteByte('.')

	la[0].deDeBruijn(builder, boundLetters, nextletter)

	// Remove our local naming
	*boundLetters = (*boundLetters)[1:]
	*nextletter = oldNextLetter
}

// String returns the lambda variable as a string
func (lv Var) String() string {
	// It is the first and only variable:
	return intToLetter(0)
}

func (lv Var) deDeBruijn(builder *strings.Builder, boundLetters *[]string, nextletter *int) {
	if int(lv) < len(*boundLetters) && (*boundLetters)[lv] != "" {
		builder.WriteString((*boundLetters)[lv])
		return
	}

	newLetter := intToLetter(*nextletter)
	*nextletter++

	for i := len(*boundLetters); i < int(lv); i++ {
		*boundLetters = append(*boundLetters, "")
	}

	*boundLetters = append(*boundLetters, newLetter)
	builder.WriteString(newLetter)
}

// String returns the lambda variable as a string
func (lf Free) String() string {
	// It is the first and only variable:
	return string(lf)
}

func (lf Free) deDeBruijn(builder *strings.Builder, boundLetters *[]string, nextletter *int) {
	builder.WriteString(string(lf))
}
