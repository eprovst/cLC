package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"

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

	return result
}
