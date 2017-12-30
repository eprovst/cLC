package main

import (
	"errors"
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
		defer close(res)
		defer close(err)

		fres, ferr := f()
		res <- fres
		err <- ferr
	}()

	select {
	case fres := <-res:
		ferr := <-err
		return fres, ferr

	case <-sig:
		return nil, errors.New("keyboard interrupt")
	}
}
