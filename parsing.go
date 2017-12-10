package LamCalc

import (
	"errors"
	"strings"
)

// ParseString turns the input into a LamTerm
func ParseString(expr string, globals map[string]LamAbst) (LamTerm, error) {
	expr = strings.TrimSpace(expr)

	if len(expr) == 0 {
		return LamExpr{}, errors.New("no expression present")
	}

	return furtherParseString(expr, map[string]int{}, globals)
}

func furtherParseString(expr string, boundVars map[string]int, globals map[string]LamAbst) (LamTerm, error) {
	var term LamTerm

	i := 0

	if expr[i] == '\\' {
		if len(expr) < 3 {
			return term, errors.New("no local variable specified in abstraction")
		}

		i++
		term = LamAbst{}

		// Create copy of boundVars where every index is one higher
		oldVars := boundVars
		boundVars = map[string]int{}

		// First increment the index of each bound variable
		for variable := range oldVars {
			boundVars[variable] = oldVars[variable] + 1
		}

		// Now get the name of the currently bound variable
		avar := ""

		for ; i < len(expr) && expr[i] != '.'; i++ {
			if expr[i] != ' ' {
				avar += string(expr[i])
			}
		}

		if i == len(expr) {
			return term, errors.New("abstraction body not started")
		}

		i++ // Skip the .

		// Add the abstraction variable to the boundvars map
		boundVars[avar] = 0

	} else {
		term = LamExpr{}
	}

	for ; i < len(expr); i++ {
		switch expr[i] {
		case '\\':
			// Start of abstraction, the rest of the expression is part of it
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
