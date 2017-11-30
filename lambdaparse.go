package LamCalc

import "reflect"

// ParseString turns the input into a LamTerm
func ParseString(expr string, boundVars map[string]int, globals map[string]LamFunc) LamTerm {
	var term LamTerm

	i := 0

	if expr[i] == 'L' {
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

		for ; expr[i] != '.'; i++ {
			fvar += string(expr[i])
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
			term = term.Append(ParseString(expr[i:], boundVars, globals))
			i = len(expr)

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

			cterm = ParseString(expr[starte:i], boundVars, globals)

			// Simplify expressions
			if reflect.TypeOf(cterm).String() == "LamCalc.LamExpr" && cterm.(LamExpr).Len() == 1 {
				cterm = cterm.(LamExpr).Index(0)
			}

			term = term.Append(cterm)

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
				term = term.Append(cindex)
			} else {
				cfnc, ok := globals[cvar]
				if ok {
					term = term.Append(cfnc)
				} else {
					panic("'" + cvar + "' not yet defined")
				}
			}
		}
	}

	// Simplify functions
	if reflect.TypeOf(term).String() == "LamCalc.LamFunc" && term.Len() == 1 {
		iterm := term.Index(0)

		if reflect.TypeOf(iterm).String() == "LamCalc.LamExpr" {
			term = LamFunc(iterm.(LamExpr))
		}
	}

	return term
}
