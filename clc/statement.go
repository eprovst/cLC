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

	case "let":
		// TODO: Is there an elegant way to stop computation?
		globals[stmnt.parameters[0].(string)] = stmnt.parameters[0].(LamCalc.LamTerm).Expand()

	case "show":
		// TODO: Is there an elegant way to stop computation?
		fmt.Print("\n" + stmnt.parameters[0].(LamCalc.LamTerm).String() + " =\n\n")
		fmt.Print("    " + stmnt.parameters[0].(LamCalc.LamTerm).Expand().String() + "\n\n")
	}
}
