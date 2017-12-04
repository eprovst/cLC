package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ElecProg/LamCalc"
)

var globals = map[string]LamCalc.LamFunc{}

func main() {
	// Prepare some basic combinators
	globals["Y"] = LamCalc.LamFunc{
		0,
		LamCalc.LamFunc{
			1,
			LamCalc.LamExpr{
				0,
				0,
			},
		},
		LamCalc.LamFunc{
			1,
			LamCalc.LamExpr{
				0,
				0,
			},
		},
	}

	globals["S"] = LamCalc.LamFunc{
		LamCalc.LamFunc{
			LamCalc.LamFunc{
				2,
				0,
				LamCalc.LamExpr{
					1,
					0,
				},
			},
		},
	}

	globals["K"] = LamCalc.LamFunc{
		LamCalc.LamFunc{
			1,
		},
	}

	globals["K*"] = LamCalc.LamFunc{
		LamCalc.LamFunc{
			0,
		},
	}

	globals["I"] = LamCalc.LamFunc{0}

	//
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])

		if err != nil {
			fmt.Println("Error:")
			fmt.Println("> " + err.Error())
			fmt.Print("Switching to interactive mode...\n\n")
		} else {
			file.Close()
			fmt.Println("File execution not yet supported...")
			os.Exit(0)
		}
	}

	commandline := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("(cLC) ")
		command, _ := commandline.ReadString('\n')
		command = strings.TrimSpace(command)

		stmnt, err := parseStatement(command)

		if err != nil {
			fmt.Println("Error: " + err.Error())
		} else {
			executeStatement(stmnt)
		}
	}
}
