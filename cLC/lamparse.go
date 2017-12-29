package main

import (
	"errors"
	"strings"

	"github.com/ElecProg/LamCalc"
)

// parseString turns the input into a LamTerm
func parseString(expr string, globals map[string]LamCalc.LamAbst) (LamCalc.LamTerm, error) {
	// Backslash is a notation for lambda
	expr = strings.Replace(expr, "\\", "λ", -1)

	return furtherParseString([]rune(expr), map[string]LamCalc.LamVar{}, globals)
}

func furtherParseString(expr []rune, boundVars map[string]LamCalc.LamVar, globals map[string]LamCalc.LamAbst) (LamCalc.LamTerm, error) {
	// Clean string
	expr = []rune(strings.TrimSpace(string(expr)))

	if len(expr) == 0 {
		return nil, errors.New("empty expression")

	} else if expr[0] == 'λ' {
		term := LamCalc.LamAbst{}
		i := 0

		if len(expr) < 2 {
			return term, errors.New("no local variable specified in abstraction")
		}

		i++

		// Create copy of boundVars where every index is one higher
		oldVars := boundVars
		boundVars = map[string]LamCalc.LamVar{}

		// First increment the index of each bound variable
		for variable := range oldVars {
			boundVars[variable] = oldVars[variable] + 1
		}

		// Now get the name of the currently bound variable
		avar := ""

		for ; i < len(expr) && expr[i] != '.'; i++ {
			avar += string(expr[i])
		}

		if !isValidVariableName(avar) {
			return term, errors.New("invalid variable name '" + avar + "'")

		} else if i >= len(expr)-1 {
			return term, errors.New("abstraction body not started")
		}

		i++ // Skip the .

		// Add the abstraction variable to the boundvars map
		boundVars[avar] = 0

		// Start of abstraction body, the rest of the expression is part of it
		part, err := furtherParseString(expr[i:], boundVars, globals)

		if err != nil {
			return nil, err
		}

		return append(term, part), nil
	}

	term := LamCalc.LamExpr{}

	for i := 0; i < len(expr); i++ {
		switch expr[i] {
		case 'λ':
			// Start of abstraction, the rest of the expression is part of it
			part, err := furtherParseString(expr[i:], boundVars, globals)
			i = len(expr)

			if err != nil {
				return nil, err
			}

			term = append(term, part)

		case '(':
			var cterm LamCalc.LamTerm

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
				return nil, err
			}

			term = append(term, cterm)

		case '\t':
			// Skip tabs

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
				term = append(term, cindex)
			} else {
				cfnc, ok := globals[cvar]
				if ok {
					term = append(term, cfnc)
				} else {
					return nil, errors.New("'" + cvar + "' not yet defined")
				}
			}
		}

		// If the LamExpr is full: encapsulate it in a new one
		if len(term) == 2 {
			term = LamCalc.LamExpr{term}
		}
	}

	// We build the LamExpr so that there is always one empty spot on the right,
	// thus we only return the first element
	return term[0], nil
}
