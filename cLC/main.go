package main

import (
	"fmt"
	"os"

	"github.com/ElecProg/LamCalc"
	"github.com/chzyer/readline"
)

var globals = map[string]LamCalc.Abst{}

func main() {
	// A warm welcome
	showInfo()

	// Limited the amount of time we wait for computation to finish
	LamCalc.MaxReductions = 5000

	// Load files
	if len(os.Args) > 1 {
		loadFiles(os.Args[1:])
		fmt.Print("Switching to interactive mode...\n\n")
	}

	commandline, _ := readline.New("(cLC) ")
	defer commandline.Close()

	for {
		command, _ := commandline.Readline()
		stmnt, err := parseStatement(command)

		if err != nil {
			printError(err)
		} else {
			executeStatement(stmnt)
		}
	}
}
