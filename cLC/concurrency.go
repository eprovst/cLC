package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"

	"github.com/ElecProg/lamcalc"
)

func concurrentReduce(term lamcalc.Term) lamcalc.Term {
	keyboardInterrupt := make(chan os.Signal, 1)
	signal.Notify(keyboardInterrupt, os.Interrupt)
	defer signal.Stop(keyboardInterrupt)

	norOut := make(chan lamcalc.Term, 1)
	stopNor := make(chan bool, 1)
	go lamcalc.ConcNorReduce(term, norOut, stopNor)

	aorOut := make(chan lamcalc.Term, 1)
	stopAor := make(chan bool, 1)
	go lamcalc.ConcAorReduce(term, aorOut, stopAor)

	var result lamcalc.Term
	select {
	case result = <-norOut:
		// Send stop signal to aor
		stopAor <- true

		// Wait for computation to stop
		<-aorOut

	case result = <-aorOut:
		// Send stop signal to not
		stopNor <- true

		// Wait for computation to stop
		<-norOut

	case <-keyboardInterrupt:
		// Remove the '^C' from the terminal:
		fmt.Print("\b\b")

		// Send stop signals
		stopNor <- true
		stopAor <- true

		// Wait for computation to stop
		<-norOut
		<-aorOut

		// Print error message
		printError(errors.New("keyboard interrupt"))
	}

	/*
	 * Okay, the following is a bit dodgy so let me explain:
	 * This program needs rediculous amounts of memory to do its calculations
	 * however once these are done the Golang runtime keeps this memory allocated
	 * in case it would need it again, immediately after the calculations. Which
	 * is here not what we want. So we force the runtime to give some memory back to
	 * the OS.
	 */
	debug.FreeOSMemory()

	return result
}
