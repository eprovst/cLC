package lambda

import (
	"errors"
	"strconv"
	"strings"
)

// TODO: Add support for free variables to the serialization format

// Deserialize turns a De Bruijn index representation
// into the internal representation
func Deserialize(inpt string) (Term, error) {
	return deserialize([]byte(inpt))
}

// MustDeserialize turns a De Bruijn index representation
// into the internal representation, panics on error
func MustDeserialize(inpt string) Term {
	result, err := deserialize([]byte(inpt))

	if err != nil {
		panic(err)
	}

	return result
}

func deserialize(inpt []byte) (Term, error) {
	var result Term

	for i := 0; i < len(inpt); i++ {
		switch inpt[i] {
		case ' ':
			// Skip spaces

		case 'l':
			body, err := deserialize(inpt[i+1:])
			i = len(inpt)

			if err != nil {
				return nil, err
			}

			if result != nil {
				result = Appl{result, Abst{body}}
			} else {
				result = Abst{body}
			}

		case '\'':
			// Skip ' (prime)
			i++

			begin := i

			for i < len(inpt) && !strings.ContainsRune(" ()", rune(inpt[i])) {
				// Aggregate the free variable
				i++
			}

			if begin >= len(inpt) || begin == i {
				return nil, errors.New("nameless free variable")
			}

			fv := Free(inpt[begin:i])

			if result != nil {
				result = Appl{result, fv}

			} else {
				result = fv
			}

			// We overshot the position of i
			i--

		case '(': // Sub expression
			nbrackets := 1
			begin := i + 1
			for nbrackets > 0 && i+1 < len(inpt) {
				i++

				switch inpt[i] {
				case '(':
					nbrackets++

				case ')':
					nbrackets--
				}
			}

			if nbrackets > 0 {
				return nil, errors.New("unbalanced brackets")
			}

			// Now i points at the closing bracket
			subexpr, err := deserialize(inpt[begin:i])

			if err != nil {
				return nil, err

			} else if result != nil {
				result = Appl{result, subexpr}

			} else {
				result = subexpr
			}

		case ')':
			return nil, errors.New("unbalanced brackets")

		default:
			begin := i
			for i+1 < len(inpt) && !strings.ContainsRune(" ()l'", rune(inpt[i+1])) {
				// Aggregate a number
				i++
			}

			idx, err := strconv.Atoi(string(inpt[begin : i+1]))
			vr := Var(idx - 1)

			if idx < 1 || err != nil {
				return nil, errors.New("invalid index: '" + string(inpt[begin:i+1]) + "'")

			} else if result != nil {
				result = Appl{result, vr}

			} else {
				result = vr
			}
		}
	}

	if result == nil {
		return nil, errors.New("empty expression")
	}

	return result, nil
}
