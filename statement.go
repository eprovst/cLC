package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/elecprog/cLC/lambda"
)

type cLCStatement struct {
	command    string
	parameters []interface{}
}

func executeStatement(stmnt cLCStatement) {
	cmd := stmnt.command

	switch cmd {
	case "none":
		// Nothing to do

	case "exit":
		os.Exit(0)

	case "help":
		showHelp()

	case "clear":
		var clearCmd *exec.Cmd

		if runtime.GOOS == "windows" {
			clearCmd = exec.Command("cmd", "/c", "cls")
		} else {
			clearCmd = exec.Command("clear")
		}

		clearCmd.Stdout = os.Stdout
		clearCmd.Run()

	case "info":
		showInfo()

	case "let", "wlet":
		var rs lambda.Term
		switch cmd {
		case "wlet":
			rs = stmnt.parameters[1].(lambda.Term).WHNF()

		default:
			rs = concurrentReduce(stmnt.parameters[1].(lambda.Term))
		}

		if rs == nil {
			// Something went wrong
			return
		}

		globals[stmnt.parameters[0].(string)] = rs

	case "free":
		// Check if all the provided variables exist
		err := existanceCheck(stmnt.parameters[0].([]string))

		if err != nil {
			printError(err)
			return
		}

		// Remove provided variables
		for _, gvar := range stmnt.parameters[0].([]string) {
			delete(globals, gvar)
		}

	case "match":
		// Check if all the provided variables exist
		err := existanceCheck(stmnt.parameters[1].([]string))

		if err != nil {
			printError(err)
			return
		}

		// Statement is perfecly fine, continue
		rs := concurrentReduce(stmnt.parameters[0].(lambda.Term))

		if rs == nil {
			// Something went wrong
			return
		}

		fmt.Print("\n" + stmnt.parameters[0].(lambda.Term).String() + " =\n")

		couldMatch := false
		for _, gvar := range stmnt.parameters[1].([]string) {
			global := globals[gvar]

			if global.AlphaEquivalent(rs) {
				fmt.Print(rs.String() + " =\n\n")
				fmt.Print("    " + gvar + "\n\n")
				couldMatch = true
				break
			}
		}

		if !couldMatch {
			fmt.Print("\n    " + rs.String() + "\n\n")
		}

	case "load":
		loadFiles(stmnt.parameters[0].([]string))

	case "show", "weak":
		var rs lambda.Term

		switch cmd {
		case "weak":
			rs = stmnt.parameters[0].(lambda.Term).WHNF()

		default:
			rs = concurrentReduce(stmnt.parameters[0].(lambda.Term))
		}

		if rs == nil {
			// Something went wrong
			return
		}

		fmt.Print("\n" + stmnt.parameters[0].(lambda.Term).String() + " =\n\n")
		fmt.Print("    " + rs.String() + "\n\n")
	}
}

// Check if all the provided variables exist
func existanceCheck(globalNames []string) error {
	nonexistant := []string{}
	for _, gvar := range globalNames {
		_, ok := globals[gvar]

		if !ok {
			nonexistant = append(nonexistant, gvar)
		}
	}

	// Build the error message if a provided variable was nonexistant
	if len(nonexistant) > 0 {
		errMsg := nonexistant[0]

		for i := 1; i < len(nonexistant)-1; i++ {
			errMsg += ", " + nonexistant[i]
		}

		if len(nonexistant) == 1 {
			errMsg += " is"

		} else {
			errMsg += " and " + nonexistant[len(nonexistant)-1] + " are"
		}

		errMsg += " not yet defined"

		return errors.New(errMsg)
	}

	return nil
}
