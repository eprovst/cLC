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

	case "match":
		nonexistant := []string{}
		// Check if all the globals exist
		for _, gvar := range stmnt.parameters[1].([]string) {
			_, ok := globals[gvar]

			if !ok {
				nonexistant = append(nonexistant, gvar)
			}
		}

		// If there are globals that don't exist...
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

			printError(errors.New(errMsg))
			return
		}

		// Statement is perfecly fine, continue
		rs := concurrentReduce(stmnt.parameters[0].(lambda.Term))

		if rs == nil {
			// Something went wrong
			return
		}

		fmt.Print("\n" + stmnt.parameters[0].(lambda.Term).String() + " =\n")

		couldFold := false
		for _, gvar := range stmnt.parameters[1].([]string) {
			global := globals[gvar]

			if global.AlphaEquivalent(rs) {
				fmt.Print(rs.String() + " =\n\n")
				fmt.Print("    " + gvar + "\n\n")
				couldFold = true
				break
			}
		}

		if !couldFold {
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
