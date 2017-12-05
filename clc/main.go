package main

import (
	"fmt"
	"os"

	"github.com/ElecProg/LamCalc"
	"github.com/chzyer/readline"
)

var globals = map[string]LamCalc.LamFunc{}

func main() {
	// Prepare some basic combinators
	// Y combinator
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

	// S, K and I calculus
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

	globals["I"] = LamCalc.LamFunc{
		0,
	}

	// Iota
	globals["i"] = LamCalc.LamFunc{
		0,
		globals["S"],
		globals["K"],
	}

	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])

		if err != nil {
			fmt.Println("Error: " + err.Error())
			fmt.Print("Switching to interactive mode...\n\n")

		} else {
			file.Close()
			fmt.Println("Error: File execution not yet supported...")
			fmt.Print("Switching to interactive mode...\n\n")
			// TODO: Support executing files.
		}
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
