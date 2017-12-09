package LamCalc

import (
	"errors"
	"strings"
)

// TODO: Improve accuracy of error messages.

// ParseString turns the input into a LamTerm
func ParseString(expr string, globals map[string]LamFunc) (LamTerm, error) {
	return furtherParseString(strings.TrimSpace(expr), map[string]int{}, globals)
}

func furtherParseString(expr string, boundVars map[string]int, globals map[string]LamFunc) (LamTerm, error) {
	var term LamTerm

	i := 0

	if expr[i] == 'L' {
		if len(expr) < 3 {
			return term, errors.New("No local variable specified in function")
		}

		i++
		term = LamFunc{}

		// Create copy of boundVars where every index is one higher
		oldVars := boundVars
		boundVars = map[string]int{}

		// First increment the index of each bound variable
		for variable := range oldVars {
			boundVars[variable] = oldVars[variable] + 1
		}

		// Now get the name of the currently bound variable
		fvar := ""

		for ; i < len(expr) && expr[i] != '.'; i++ {
			if expr[i] != ' ' {
				fvar += string(expr[i])
			}
		}

		if i == len(expr) {
			return term, errors.New("Function body not started")
		}

		i++ // Skip the .

		// Add the function variable to the boundvars map
		boundVars[fvar] = 0

	} else {
		term = LamExpr{}
	}

	for ; i < len(expr); i++ {
		switch expr[i] {
		case 'L':
			// Start of function, the rest of the expression is part of it
			part, err := furtherParseString(expr[i:], boundVars, globals)
			i = len(expr)

			if err != nil {
				return term, err
			}

			term = term.append(part)

		case '(':
			var cterm interface{}

			i++
			starte := i

			nBrack := 0
			for ; i < len(expr); i++ {
				if expr[i] == ')' {
					if nBrack == 0 {
						break
					} else {
						nBrack--
					}
				} else if expr[i] == '(' {
					nBrack++
				}
			} // After this loop i points at the closing bracket

			cterm, err := furtherParseString(expr[starte:i], boundVars, globals)

			if err != nil {
				return term, err
			}

			term = term.append(cterm)

		case ' ':
			// Skip spaces

		default:
			// A variable
			cvar := ""

			for ; i < len(expr); i++ {
				if expr[i] == '(' {
					// End of var, take another look at the character later on
					i--
					break
				} else if expr[i] == ' ' {
					// Space end of var
					break
				} else {
					cvar += string(expr[i])
				}
			}

			cindex, ok := boundVars[cvar]

			if ok {
				term = term.append(cindex)
			} else {
				cfnc, ok := globals[cvar]
				if ok {
					term = term.append(cfnc)
				} else {
					return term, errors.New("'" + cvar + "' not yet defined")
				}
			}
		}
	}

	return term, nil
}
