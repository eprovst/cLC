package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
)

// Here we implement a system to call a function and make it stoppable on ^+C
// The function isn't really stopped, rather it's result is discarded, and we
// don't wait for it anymore either.

func stoppable(f func() (interface{}, error)) (interface{}, error) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	defer signal.Stop(sig)

	res := make(chan interface{})
	err := make(chan error)

	go func() {
		fres, ferr := f()
		res <- fres
		err <- ferr
	}()

	select {
	case fres := <-res:
		ferr := <-err
		return fres, ferr

	case <-sig:
		// Remove the '^C' from the terminal:
		fmt.Print("\b\b")

		return nil, errors.New("keyboard interrupt")
	}
}
