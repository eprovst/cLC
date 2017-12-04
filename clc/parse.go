package main

import (
	"strings"

	"github.com/ElecProg/LamCalc"
)

func parseStatement(stmnt string) (cLCStatement, error) {
	// Some clean up
	stmnt = strings.TrimSpace(stmnt)

	// TODO: Add let
	switch strings.Fields(stmnt)[0] {
	case "exit":
		return cLCStatement{command: "exit"}, nil

	case "clear":
		return cLCStatement{command: "clear"}, nil

	case "let":
		// TODO: Make more robust
		stmnt = strings.TrimPrefix(stmnt, "let")
		varname := strings.TrimSpace(strings.TrimSuffix(strings.SplitAfter(stmnt, "=")[0], "="))
		expression, err := LamCalc.ParseString(strings.TrimSpace(strings.SplitAfter(stmnt, "=")[1]), map[string]int{}, globals)

		if err != nil {
			return cLCStatement{}, err
		}

		return cLCStatement{
			command:    "let",
			parameters: []interface{}{varname, expression},
		}, nil

	default:
		expression, err := LamCalc.ParseString(stmnt, map[string]int{}, globals)

		if err != nil {
			return cLCStatement{}, err
		}

		return cLCStatement{
			command:    "show",
			parameters: []interface{}{expression},
		}, nil
	}
}
