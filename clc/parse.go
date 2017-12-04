package main

import (
	"strings"

	"github.com/ElecProg/LamCalc"
)

func parseStatement(stmnt string) (cLCStatement, error) {
	// TODO: Add let
	switch strings.Fields(stmnt)[0] {
	case "exit":
		return cLCStatement{command: "exit"}, nil

	case "clear":
		return cLCStatement{command: "clear"}, nil

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
