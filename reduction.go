package LamCalc

import (
	"errors"
	"reflect"
)

// MaxReductions determines the maximum amount of expansions before we give up
var MaxReductions = 10000

// reduceOnce reduces a lambda expression once
func (lx LamExpr) reduceOnce() LamTerm {
	nw := LamExpr{}

	if len(lx) >= 2 && reflect.TypeOf(lx[0]).String() == "LamCalc.LamFunc" {
		nw = append(nw, lx[0].(LamFunc).betaReduce(lx[1]))

		if len(lx) > 2 {
			nw = append(nw, lx[2:]...)
		}

		return nw
	}

	for _, term := range lx {
		switch term := term.(type) {
		case int:
			nw = append(nw, term)

		case LamTerm:
			nw = append(nw, term.reduceOnce())
		}
	}

	return nw
}

func (lf LamFunc) reduceOnce() LamTerm {
	return LamFunc{LamExpr(lf).reduceOnce()}
}

// Reduce reduces a lambda expression
func (lx LamExpr) Reduce() (LamFunc, error) {
	ls := lx.simplify()
	nw := ls.reduceOnce().simplify()

	for c := 0; !nw.alphaEquivalent(ls); c++ {
		if c == MaxReductions {
			return LamFunc{}, errors.New("exeeded maximum amount of reductions")
		}

		ls = nw
		nw = nw.reduceOnce().simplify()
	}

	return nw.(LamFunc).etaReduce().simplify().(LamFunc), nil
}

// Reduce reduces a lambda function
func (lf LamFunc) Reduce() (LamFunc, error) {
	ls := lf.simplify()
	nw := ls.reduceOnce().simplify()

	for c := 0; !nw.alphaEquivalent(ls); c++ {
		if c == MaxReductions {
			return LamFunc{}, errors.New("exeeded maximum amount of reductions")
		}

		ls = nw
		nw = nw.reduceOnce().simplify()
	}

	return nw.(LamFunc).etaReduce().simplify().(LamFunc), nil
}

// WHNFReduce reduces the expression till a weak head normal form is found (eta reduction isn't tried)
func (lx LamExpr) WHNFReduce() (LamFunc, error) {
	nw := lx.reduceOnce().simplify()

	for c := 0; reflect.TypeOf(nw).String() != "LamCalc.LamFunc"; c++ {
		if c == MaxReductions {
			return LamFunc{}, errors.New("exeeded maximum amount of reductions")
		}

		nw = nw.reduceOnce().simplify()
	}

	return nw.(LamFunc), nil
}

// WHNFReduce reduces the function till a weak head normal form is found (eta reduction isn't tried)
// ie. doesn't do anything...
func (lf LamFunc) WHNFReduce() (LamFunc, error) {
	return lf, nil
}
