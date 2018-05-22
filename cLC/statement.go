package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/ElecProg/lamcalc"
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

	case "let", "alet", "wlet":
		var rs interface{}
		var err error

		switch cmd {
		case "alet":
			rs, err = stoppable(func() (interface{}, error) {
				return stmnt.parameters[1].(lamcalc.Term).AorReduce()
			})

		case "wlet":
			rs, err = stmnt.parameters[1].(lamcalc.Term), error(nil)

		default:
			rs, err = stoppable(func() (interface{}, error) {
				return stmnt.parameters[1].(lamcalc.Term).NorReduce()
			})
		}

		if err != nil {
			printError(err)
			return
		}

		lx := rs.(lamcalc.Term)

		// Make sure it's a function
		la := lx.WHNF()

		globals[stmnt.parameters[0].(string)] = la

	case "fold":
		rs, err := stoppable(func() (interface{}, error) {
			return stmnt.parameters[1].(lamcalc.Term).NorReduce()
		})

		if err != nil {
			printError(err)
			return
		}

		expression := rs.(lamcalc.Term)

		fmt.Print("\n" + stmnt.parameters[0].(lamcalc.Term).String() + " =\n")

		couldFold := false
		for _, gvar := range stmnt.parameters[1].([]string) {
			if globals[gvar].AlphaEquivalent(expression) {
				fmt.Print(expression.String() + " =\n\n")
				fmt.Print("    " + gvar + "\n\n")
				couldFold = true
				break
			}
		}

		if !couldFold {
			fmt.Print("\n    " + expression.String() + "\n\n")
		}

	case "load":
		loadFiles(stmnt.parameters[0].([]string))

	case "show", "apor", "weak":
		var rs interface{}
		var err error

		switch cmd {
		case "apor":
			rs, err = stoppable(func() (interface{}, error) {
				return stmnt.parameters[0].(lamcalc.Term).AorReduce()
			})

		case "weak":
			rs, err = stmnt.parameters[0].(lamcalc.Term).WHNF(), error(nil)

		default:
			rs, err = stoppable(func() (interface{}, error) {
				return stmnt.parameters[0].(lamcalc.Term).NorReduce()
			})
		}

		if err != nil {
			printError(err)
			return
		}

		expression := rs.(lamcalc.Term)

		fmt.Print("\n" + stmnt.parameters[0].(lamcalc.Term).String() + " =\n\n")
		fmt.Print("    " + expression.String() + "\n\n")
	}
}
