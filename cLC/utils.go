package main

import (
	"bufio"
	"fmt"
	"os"
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

cLamCalc
--------

commandline Lambda Calculator

Copyright (c) 2017 Evert Provoost.
All Rights Reserved.

`)
}
