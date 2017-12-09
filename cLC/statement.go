package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/ElecProg/LamCalc"
)

type cLCStatement struct {
	command    string
	parameters []interface{}
}

func executeStatement(stmnt cLCStatement) {
	switch stmnt.command {
	case "none":
		// Nothing to do
	case "exit":
		os.Exit(0)

	case "clear":
		var clearCmd *exec.Cmd

		if runtime.GOOS == "windows" {
			clearCmd = exec.Command("cls")
		} else {
			clearCmd = exec.Command("clear")
		}

		clearCmd.Stdout = os.Stdout
		clearCmd.Run()

	case "info":
		showInfo()

	case "let":
		// TODO: Is there an elegant way to stop computation?
		lf := stmnt.parameters[1].(LamCalc.LamTerm).Reduce()
		globals[stmnt.parameters[0].(string)] = lf

	case "fold":
		// TODO: Is there an elegant way to stop computation?
		fmt.Print("\n" + stmnt.parameters[0].(LamCalc.LamTerm).String() + " =\n")

		expression := stmnt.parameters[0].(LamCalc.LamTerm).Reduce()

		couldFold := false
		for _, term := range stmnt.parameters[1].([]string) {
			if globals[term].Equivalent(expression) {
				fmt.Print(expression.String() + " =\n\n")
				fmt.Print("    " + term + "\n\n")
				couldFold = true
				break
			}
		}

		if !couldFold {
			fmt.Print("\n    " + expression.String() + "\n\n")
		}

	case "load":
		loadFiles(stmnt.parameters[0].([]string))

	case "show":
		// TODO: Is there an elegant way to stop computation?
		fmt.Print("\n" + stmnt.parameters[0].(LamCalc.LamTerm).String() + " =\n\n")
		fmt.Print("    " + stmnt.parameters[0].(LamCalc.LamTerm).Reduce().String() + "\n\n")
	}
}
