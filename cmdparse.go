package main

import (
	"errors"
	"strings"
)

func parseStatement(stmnt string) (cLCStatement, error) {
	// Remove comments
	stmnt = strings.SplitAfter(stmnt, "--")[0]
	stmnt = strings.TrimSuffix(stmnt, "--")
	stmnt = strings.TrimSpace(stmnt)

	if len(stmnt) == 0 {
		// Empty line
		return cLCStatement{command: "none"}, nil
	}

	cmd := strings.Fields(stmnt)[0]

	switch cmd {
	case "exit", "clear", "info", "help":
		return cLCStatement{command: cmd}, nil

	case "let", "wlet":
		stmnt = strings.TrimPrefix(stmnt, cmd)
		splitStmnt := strings.SplitN(stmnt, "=", 2)

		if len(splitStmnt) < 2 {
			return cLCStatement{}, errors.New("no expression in " + cmd + " operation")
		}

		varname := strings.TrimSpace(splitStmnt[0])
		// \ should always become λ
		varname = strings.Replace(varname, "\\", "λ", -1)

		if !isValidVariableName(varname) {
			return cLCStatement{}, errors.New("invalid variable name '" + varname + "' in " + cmd + " operation")
		}

		expression, err := parseString(splitStmnt[1], globals)

		if err != nil {
			return cLCStatement{}, err
		}

		return cLCStatement{
			command:    cmd,
			parameters: []interface{}{varname, expression},
		}, nil

	case "free":
		stmnt = strings.TrimPrefix(stmnt, "free")
		vars := strings.Fields(stmnt)

		if len(vars) == 0 {
			return cLCStatement{}, errors.New("no targets in free operation")
		}

		return cLCStatement{
			command:    "free",
			parameters: []interface{}{vars},
		}, nil

	case "match":
		stmnt = strings.TrimPrefix(stmnt, "match")
		splitStmnt := strings.SplitN(stmnt, "with", 2)

		if len(splitStmnt) < 2 {
			return cLCStatement{}, errors.New("no targets in match operation")
		}

		expression, err := parseString(splitStmnt[0], globals)
		vars := strings.Fields(splitStmnt[1])

		if len(vars) == 0 {
			return cLCStatement{}, errors.New("no targets in match operation")

		} else if err != nil {
			return cLCStatement{}, err
		}

		return cLCStatement{
			command:    "match",
			parameters: []interface{}{expression, vars},
		}, nil

	case "load":
		fields := strings.Fields(stmnt)

		if len(fields) > 1 {
			return cLCStatement{
				command:    "load",
				parameters: []interface{}{fields[1:]},
			}, nil
		}

		return cLCStatement{}, errors.New("no files listed to load")

	case "weak":
		stmnt = strings.TrimPrefix(stmnt, cmd)
		expression, err := parseString(stmnt, globals)

		if err != nil {
			return cLCStatement{}, err
		}

		return cLCStatement{
			command:    cmd,
			parameters: []interface{}{expression},
		}, nil

	default:
		expression, err := parseString(stmnt, globals)

		if err != nil {
			return cLCStatement{}, err
		}

		return cLCStatement{
			command:    "show",
			parameters: []interface{}{expression},
		}, nil
	}
}
