package LamCalc

import (
	"reflect"
)

// reduceOnce expands a lambda expression once
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

	return nw.(LamFunc)
}

// Reduce reduces a lambda function
func (lf LamFunc) Reduce() LamFunc {
	ls := lf.simplify()
	nw := ls.reduceOnce().simplify()

	for !nw.alphaEquivalent(ls) {
		ls = nw
		nw = nw.reduceOnce().simplify()
	}

	return nw.(LamFunc)
}
