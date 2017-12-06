package main

import (
	"strings"

	"github.com/ElecProg/LamCalc"
)

func parseStatement(stmnt string) (cLCStatement, error) {
	// Some clean up
	stmnt = strings.TrimSpace(stmnt)

	if len(stmnt) == 0 || strings.HasPrefix(stmnt, "--") {
		// Empty line or comment
		return cLCStatement{command: "none"}, nil
	}

	// TODO: Add let
	switch strings.Fields(stmnt)[0] {
	case "exit":
		return cLCStatement{command: "exit"}, nil

	case "clear":
		return cLCStatement{command: "clear"}, nil

	case "info":
		return cLCStatement{command: "info"}, nil

	case "let":
		// TODO: Make more robust
		stmnt = strings.TrimPrefix(stmnt, "let")
		varname := strings.TrimSpace(strings.TrimSuffix(strings.SplitAfter(stmnt, "=")[0], "="))
		expression, err := LamCalc.ParseString(strings.SplitAfter(stmnt, "=")[1], globals)

		if err != nil {
			return cLCStatement{}, err
		}

		return cLCStatement{
			command:    "let",
			parameters: []interface{}{varname, expression},
		}, nil

	case "fold":
		// TODO: Make more robust
		stmnt = strings.TrimPrefix(stmnt, "fold")
		expression, err := LamCalc.ParseString(strings.TrimSuffix(strings.SplitAfter(stmnt, "into")[0], "into"), globals)
		vars := strings.Fields(strings.SplitAfter(stmnt, "into")[1])

		if err != nil {
			return cLCStatement{}, err
		}

		return cLCStatement{
			command:    "fold",
			parameters: []interface{}{expression, vars},
		}, nil

	default:
		expression, err := LamCalc.ParseString(stmnt, globals)

		if err != nil {
			return cLCStatement{}, err
		}

		return cLCStatement{
			command:    "show",
			parameters: []interface{}{expression},
		}, nil
	}
}
