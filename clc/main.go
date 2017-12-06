package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ElecProg/LamCalc"
	"github.com/chzyer/readline"
)

var globals = map[string]LamCalc.LamFunc{}

func main() {
	// Load files
	if len(os.Args) > 1 {
		for _, filePath := range os.Args[1:] {
			file, err := os.Open(filePath)

			if err != nil {
				fmt.Println("Error: " + err.Error())

			} else {
				fileScanner := bufio.NewScanner(file)

				for fileScanner.Scan() {
					command := fileScanner.Text()
					stmnt, err := parseStatement(command)

					if err != nil {
						fmt.Println("Error: " + err.Error())
					} else {
						executeStatement(stmnt)
					}
				}

				file.Close()
				fmt.Println("Done loading '" + filePath + "'.")
			}
		}

		fmt.Print("Switching to interactive mode...\n\n")
	}

	commandline, _ := readline.New("(cLC) ")
	defer commandline.Close()

	// Show info
	executeStatement(cLCStatement{command: "info"})

	for {
		command, _ := commandline.Readline()
		stmnt, err := parseStatement(command)

		if err != nil {
			fmt.Println("Error: " + err.Error())
		} else {
			executeStatement(stmnt)
		}
	}
}
