package LamCalc

import (
	"reflect"
)

// Substitute replaces index by sub
func (lx LamExpr) substitute(index int, sub LamTerm) LamExpr {
	nw := LamExpr{}

	for _, term := range lx {
		switch term := term.(type) {
		case int:
			if term == index {
				nw = append(nw, sub)

			} else {
				nw = append(nw, term)
			}

		case LamExpr:
			nw = append(nw, term.substitute(index, sub))

		case LamFunc:
			nw = append(nw, term.substitute(index+1, sub))
		}
	}

	return nw
}

// Substitute replaces index by sub
func (lf LamFunc) substitute(index int, sub LamTerm) LamFunc {
	return LamFunc(LamExpr(lf).substitute(index, sub))
}

// Insert replaces index 0 by sub and returns a LamExpr
func (lf LamFunc) insert(sub LamTerm) LamExpr {
	return LamExpr(lf).substitute(0, sub)
}

// Expand expands a lambda function (which means doing, mostly, nothing)
func (lf LamFunc) Expand() LamFunc {
	return lf.simplify()
}

// ExpandOnce expands a lambda expression once
func (lx LamExpr) expandOnce() LamTerm {
	nw := LamExpr{}

	if len(lx) == 1 {
		return lx[0].(LamTerm)

	} else if reflect.TypeOf(lx[0]).String() == "LamCalc.LamFunc" &&
		reflect.TypeOf(lx[1]).String() == "LamCalc.LamFunc" {

		nw = append(nw, lx[0].(LamFunc).insert(lx[1].(LamFunc)))

		if len(lx) > 2 {
			nw = append(nw, lx[2:]...)
		}

		return nw
	}

	for _, term := range lx {
		switch term := term.(type) {
		case int:
			nw = append(nw, term)

		case LamFunc:
			nw = append(nw, term)

		case LamExpr:
			nw = append(nw, term.expandOnce())
		}
	}

	return nw
}

// Expand expands a lambda expression
func (lx LamExpr) Expand() LamFunc {
	nw := lx.simplify().expandOnce()

	for reflect.TypeOf(nw).String() != "LamCalc.LamFunc" {
		nw = nw.(LamExpr).expandOnce()
	}

	return nw.(LamFunc)
}

// Simplify (tries) to remove unnecessary brackets
func (lx LamExpr) simplify() LamExpr {
	if reflect.TypeOf(lx[0]).String() == "LamCalc.LamExpr" {
		res := lx[0].(LamExpr)

		if len(lx) > 1 {
			res = append(res, lx[1:]...)
		}

		return res.simplify()
	}

	res := LamExpr{}

	for _, term := range lx {
		switch term := term.(type) {
		case LamExpr:
			simpl := term.simplify()

			if len(simpl) == 1 {
				res = append(res, simpl[0])
			} else {
				res = append(res, simpl)
			}

		case LamFunc:
			res = append(res, term.simplify())

		default:
			res = append(res, term)
		}
	}

	return res
}

// Simplify (tries) to remove unnecessary brackets
func (lf LamFunc) simplify() LamFunc {
	return LamFunc(LamExpr(lf).simplify())
}
