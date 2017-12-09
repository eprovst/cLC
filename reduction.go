package LamCalc

import (
	"reflect"
)

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
func (lx LamExpr) Reduce() LamFunc {
	ls := lx.simplify()
	nw := ls.reduceOnce().simplify()

	for !nw.alphaEquivalent(ls) {
		ls = nw
		nw = nw.reduceOnce().simplify()
	}

	return nw.(LamFunc).etaReduce().simplify().(LamFunc)
}

// Reduce reduces a lambda function
func (lf LamFunc) Reduce() LamFunc {
	ls := lf.simplify()
	nw := ls.reduceOnce().simplify()

	for !nw.alphaEquivalent(ls) {
		ls = nw
		nw = nw.reduceOnce().simplify()
	}

	return nw.(LamFunc).etaReduce().simplify().(LamFunc)
}

// HNFReduce reduces the expression till a head normal form is found (eta reduction isn't tried)
func (lx LamExpr) HNFReduce() LamFunc {
	nw := lx.reduceOnce().simplify()

	for reflect.TypeOf(nw).String() != "LamCalc.LamFunc" {
		nw = nw.reduceOnce().simplify()
	}

	return nw.(LamFunc)
}

// HNFReduce reduces the function till a head normal form is found (eta reduction isn't tried)
// ie. doesn't do anything...
func (lf LamFunc) HNFReduce() LamFunc {
	return lf
}
