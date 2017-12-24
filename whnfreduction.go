package LamCalc

import (
	"errors"
	"reflect"
)

// WHNFReduce reduces the expression till a weak head normal form is found (eta reduction isn't tried)
func (lx LamExpr) WHNFReduce() (LamAbst, error) {
	nw := lx.Simplify()

	for c := 0; reflect.TypeOf(nw).String() != "LamCalc.LamAbst"; c++ {
		if c == MaxReductions {
			return LamAbst{}, errors.New("exeeded maximum amount of reductions")
		}

		nw = nw.reduceOnce().Simplify()
	}

	return nw.(LamAbst), nil
}

// WHNFReduce reduces the abstraction till a weak head normal form is found (eta reduction isn't tried)
// ie. doesn't do anything
func (la LamAbst) WHNFReduce() (LamAbst, error) {
	return la, nil
}

// WHNFReduce won't work on a free variable
func (lv LamVar) WHNFReduce() (LamAbst, error) {
	return nil, errors.New("can't reduce free variable to WHNF")
}
