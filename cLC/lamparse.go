package main

import (
	"errors"
	"strings"

	"github.com/ElecProg/LamCalc"
)

// parseString turns the input into a LamTerm
func parseString(expr string, globals map[string]LamCalc.LamAbst) (LamCalc.LamTerm, error) {
	expr = strings.TrimSpace(expr)

	// Backslash is a notation for lambda
	expr = strings.Replace(expr, "\\", "λ", -1)

	if len(expr) == 0 {
		return LamCalc.LamExpr{}, errors.New("no expression present")
	}

	return furtherParseString([]rune(expr), map[string]LamCalc.LamVar{}, globals)
}

func furtherParseString(expr []rune, boundVars map[string]LamCalc.LamVar, globals map[string]LamCalc.LamAbst) (LamCalc.LamTerm, error) {
	var term LamCalc.LamTerm

	i := 0

	if expr[i] == 'λ' {
		if len(expr) < 2 {
			return term, errors.New("no local variable specified in abstraction")
		}

		i++
		term = LamCalc.LamAbst{}

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

	} else {
		term = LamCalc.LamExpr{}
	}

	for ; i < len(expr); i++ {
		switch expr[i] {
		case 'λ':
			// Start of abstraction, the rest of the expression is part of it
			part, err := furtherParseString(expr[i:], boundVars, globals)
			i = len(expr)

			if err != nil {
				return term, err
			}

			term = appendToLT(term, part)

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
				return term, err
			}

			term = appendToLT(term, cterm)

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
				term = appendToLT(term, cindex)
			} else {
				cfnc, ok := globals[cvar]
				if ok {
					term = appendToLT(term, cfnc)
				} else {
					return term, errors.New("'" + cvar + "' not yet defined")
				}
			}
		}
	}

	return term, nil
}

// Utility to add to a LamCalc.LamTerm
func appendToLT(original LamCalc.LamTerm, toAdd ...LamCalc.LamTerm) LamCalc.LamTerm {
	switch original.(type) {
	case LamCalc.LamAbst:
		return append(original.(LamCalc.LamAbst), toAdd...)

	default:
		return append(original.(LamCalc.LamExpr), toAdd...)
	}
}
