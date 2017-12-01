package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ElecProg/LamCalc"
)

func main() {
	// Prepare some basic combinators
	Y := LamCalc.LamFunc{
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

	S := LamCalc.LamFunc{
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

	K := LamCalc.LamFunc{
		LamCalc.LamFunc{
			1,
		},
	}

	Ks := LamCalc.LamFunc{
		LamCalc.LamFunc{
			0,
		},
	}

	I := LamCalc.LamFunc{0}

	globals := map[string]LamCalc.LamFunc{
		"Y":  Y,
		"S":  S,
		"K":  K,
		"K*": Ks,
		"I":  I,
	}

	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])

		if err != nil {
			fmt.Println("Error:")
			fmt.Println("> " + err.Error())
			fmt.Print("Switching to interactive mode...\n\n")
		} else {
			file.Close()
			os.Exit(0)
		}
	}

	commandline := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("(cLC) ")
		command, _ := commandline.ReadString('\n')
		command = strings.TrimSpace(command)

		// Commandline commands
		// TODO: Extend functionality
		switch command {
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

		default:
			// TODO: Is there an elemgant way to stop computation?
			lx, err := LamCalc.ParseString(command, map[string]int{}, globals)

			if err != nil {
				fmt.Println("Error: " + err.Error())
			} else {
				fmt.Print("\n" + lx.String() + " =\n\n")
				fmt.Print("    " + lx.Expand().String() + "\n\n")
			}
		}
	}
}
