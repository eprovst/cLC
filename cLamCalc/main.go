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
	// Prepare some sort of demo
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
			fmt.Print("\nY =\n\n")
			fmt.Print("    " + Y.String() + "\n\n")
		}
	}
}
