package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func printError(err error) {
	fmt.Println("Error: " + err.Error())
}

func loadFiles(paths []string) {
	for _, filePath := range paths {
		file, err := os.Open(filePath)

		if err != nil {
			printError(err)

		} else {
			fileScanner := bufio.NewScanner(file)

			for fileScanner.Scan() {
				command := fileScanner.Text()
				stmnt, err := parseStatement(command)

				if err != nil {
					printError(err)
				} else {
					executeStatement(stmnt)
				}
			}

			file.Close()
			fmt.Println("Done loading '" + filePath + "'.")
		}
	}
}

func showInfo() {
	fmt.Print(`       _      ___
   __   \    /
  /     /\  (
  \__  /  \  \___

cLamCalc v1.2.0
---------------

commandline Lambda Calculator

Copyright (c) 2017-2018 Evert Provoost.
All Rights Reserved.

`)
}

func showHelp() {
	fmt.Print(`
Help:
-----

For full details: visit the project's wiki.

Availabe commands:

<lambda expression>
→ Normal order and applicative order expansion are tried for the expression, if there's a result it will be shown.

let <new global> = <lambda expression>
→ If the expansion can be fully reduced sets the global equal to that reduced form.

fold <lambda expression> into <global1> <global2> <...>
→ Tries to fully expand the expression and then shows the first listed global which is equivalent to that reduction.

weak <lambda expression>
→ Transforms the expression to a weak head normal form, then shows the result. Useful for expressions which wouldn't terminate reducing otherwise.

wlet <new global> = <lambda expression>
→ Equivalent to let but only transforms the expression to a weak head normal form.

<command> -- <comment>
→ Everything after -- is ignored.

help
→ Shows help for the cLC.

clear
→ Clears the terminal.

info
→ Shows information about the cLC.

exit
→ Closes the cLC.

`)
}

func isValidVariableName(varname string) bool {
	return !strings.HasPrefix(varname, "\\") && !strings.HasPrefix(varname, "λ") && !strings.ContainsAny(varname, " \t")
}
