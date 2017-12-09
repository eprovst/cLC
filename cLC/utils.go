package main

import (
	"bufio"
	"fmt"
	"os"
)

func loadFiles(paths []string) {
	for _, filePath := range paths {
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
