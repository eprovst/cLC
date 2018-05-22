package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/ElecProg/lamcalc"
)

// Here we implement a system to call a function and make it stoppable on ^+C
// The function isn't really stopped, rather it's result is discarded, and we
// don't wait for it anymore either.

func concurrentReduce(term lamcalc.Term) (lamcalc.Term, error) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	defer signal.Stop(sig)

	resn := make(chan lamcalc.Term)
	errn := make(chan error)

	go func() {
		fres, ferr := term.NorReduce()
		resn <- fres
		errn <- ferr
	}()

	resa := make(chan lamcalc.Term)
	erra := make(chan error)

	go func() {
		fres, ferr := term.AorReduce()
		resa <- fres
		erra <- ferr
	}()

	var err error

	for errs := 0; errs < 2; errs++ {
		select {
		case fres := <-resn:
			ferr := <-errn

			if ferr == nil {
				return fres, nil
			}

			err = ferr

		case fres := <-resa:
			ferr := <-erra

			if ferr == nil {
				return fres, nil
			}

			err = ferr

		case <-sig:
			// Remove the '^C' from the terminal:
			fmt.Print("\b\b")

			return nil, errors.New("keyboard interrupt")
		}
	}

	// Both returned errors:
	return nil, err
}
