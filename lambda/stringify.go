package lambda

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

func getFreeVariables(t Term) []Free {
	switch t := t.(type) {
	case Free:
		return []Free{t}

	case Appl:
		return append(getFreeVariables(t[0]), getFreeVariables(t[1])...)

	case Abst:
		return getFreeVariables(t[0])

	default:
		return []Free{}
	}
}

func in(t Free, f *[]Free) bool {
	for _, v := range *f {
		if t == v {
			return true
		}
	}

	return false
}

// String returns the Lambda Expression as a string
func (lx Appl) String() string {
	builder := strings.Builder{}
	freeVars := getFreeVariables(lx)
	lx.deDeBruijn(&builder, new([]string), new(int), &freeVars)

	return builder.String()
}

func (lx Appl) deDeBruijn(builder *strings.Builder, boundLetters *[]string, nextletter *int, freeVars *[]Free) {
	for i, part := range lx {
		switch part := part.(type) {
		case Abst:
			builder.WriteByte('(')
			part.deDeBruijn(builder, boundLetters, nextletter, freeVars)
			builder.WriteByte(')')

		case Appl:
			if i == 0 {
				part.deDeBruijn(builder, boundLetters, nextletter, freeVars)

			} else {
				builder.WriteByte('(')
				part.deDeBruijn(builder, boundLetters, nextletter, freeVars)
				builder.WriteByte(')')
			}

		default:
			part.deDeBruijn(builder, boundLetters, nextletter, freeVars)
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
	freeVars := getFreeVariables(la)
	la.deDeBruijn(&builder, new([]string), new(int), &freeVars)

	return builder.String()
}

func (la Abst) deDeBruijn(builder *strings.Builder, boundLetters *[]string, nextletter *int, freeVars *[]Free) {
	// Remember at which character we were
	oldNextLetter := *nextletter

	// Search for a free name
	newLetter := intToLetter(*nextletter)
	*nextletter++

	for in(Free(newLetter), freeVars) {
		newLetter = intToLetter(*nextletter)
		*nextletter++
	}

	// Add the localy bound letter to the list
	*boundLetters = append([]string{newLetter}, *boundLetters...)

	builder.WriteRune('λ')
	builder.WriteString(newLetter)
	builder.WriteByte('.')

	la[0].deDeBruijn(builder, boundLetters, nextletter, freeVars)

	// Remove our local naming
	*boundLetters = (*boundLetters)[1:]
	*nextletter = oldNextLetter
}

// String returns the lambda variable as a string
func (lv Var) String() string {
	// It is the first and only variable:
	return intToLetter(0)
}

func (lv Var) deDeBruijn(builder *strings.Builder, boundLetters *[]string, nextletter *int, freeVars *[]Free) {
	if int(lv) < len(*boundLetters) && (*boundLetters)[lv] != "" {
		builder.WriteString((*boundLetters)[lv])
		return
	}

	// Search for a free name
	newLetter := intToLetter(*nextletter)
	*nextletter++

	for in(Free(newLetter), freeVars) {
		newLetter = intToLetter(*nextletter)
		*nextletter++
	}

	// Add it to the boundletters table
	for i := len(*boundLetters); i < int(lv); i++ {
		*boundLetters = append(*boundLetters, "")
	}

	*boundLetters = append(*boundLetters, newLetter)

	// And finaly write it to output
	builder.WriteString(newLetter)
}

// String returns the lambda variable as a string
func (lf Free) String() string {
	// It is the first and only variable:
	return string(lf)
}

func (lf Free) deDeBruijn(builder *strings.Builder, boundLetters *[]string, nextletter *int, freeVars *[]Free) {
	builder.WriteString(string(lf))
}
