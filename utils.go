package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func printError(err error) {
	fmt.Println("Error: " + err.Error())
}

func loadFiles(paths []string) {
	for _, filePath := range paths {
		// Get the current working directory
		currentPath, err := os.Getwd()

		if err != nil {
			printError(err)
			continue
		}

		// Open the file
		file, err := os.Open(filePath)

		if err != nil {
			printError(err)
			continue
		}

		// Go to the file's directory
		err = os.Chdir(filepath.Dir(filePath))
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

		// Go back to the original path
		err = os.Chdir(currentPath)

		if err != nil {
			printError(err)
		}
	}
}

func showInfo() {
	fmt.Print(`               _      ___
           __   \    /
          /     /\  (
          \__  /  \  \___

commandline Lambda Calculator v1.4.0
------------------------------------

Copyright (c) 2017-2018 Evert Provoost.
Some rights reserved.

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

free <global1> <global2> <...>
→ Unbinds the global(s) and thus makes it a free variable.

match <lambda expression> with <global1> <global2> <...>
→ Tries to fully expand the expression and then shows the first listed global variable which is equivalent to that reduction.

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
