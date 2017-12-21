package LamCalc

import (
	"errors"
	"reflect"
)

// MaxReductions determines the maximum amount of expansions before we give up
// use a negative value to have no limit (use with care...)
var MaxReductions = 10000

// reduceOnce reduces a lambda expression once
func (lx LamExpr) reduceOnce() LamTerm {
	nw := LamExpr{}

	if len(lx) >= 2 && reflect.TypeOf(lx[0]).String() == "LamCalc.LamAbst" {
		nw = append(nw, lx[0].(LamAbst).betaReduce(lx[1]))

		if len(lx) > 2 {
			nw = append(nw, lx[2:]...)
		}

		return nw
	}

	for _, term := range lx {
		switch term := term.(type) {
		case LamVar:
			nw = append(nw, term)

		case LamTerm:
			nw = append(nw, term.reduceOnce())
		}
	}

	return nw
}

func (la LamAbst) reduceOnce() LamTerm {
	return LamAbst{LamExpr(la).reduceOnce()}
}

func (lv LamVar) reduceOnce() LamTerm {
	return lv
}

// Reduce reduces a lambda expression using normal order
func (lx LamExpr) Reduce() (LamTerm, error) {
	ls := lx.Simplify()
	nw := ls.reduceOnce().Simplify()

	for c := 0; !nw.alphaEquivalent(ls); c++ {
		if c == MaxReductions {
			return LamAbst{}, errors.New("exeeded maximum amount of reductions")
		}

		ls = nw
		nw = nw.reduceOnce().Simplify()
	}

	return nw.etaReduce().Simplify(), nil
}

// Reduce reduces a lambda abstraction using normal order
func (la LamAbst) Reduce() (LamTerm, error) {
	ls := la.Simplify()
	nw := ls.reduceOnce().Simplify()

	for c := 0; !nw.alphaEquivalent(ls); c++ {
		if c == MaxReductions {
			return LamAbst{}, errors.New("exeeded maximum amount of reductions")
		}

		ls = nw
		nw = nw.reduceOnce().Simplify()
	}

	return nw.etaReduce().Simplify(), nil
}

// Reduce won't work on a free variable
func (lv LamVar) Reduce() (LamTerm, error) {
	return nil, errors.New("can't reduce free variable")
}

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
