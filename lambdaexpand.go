package LamCalc

import (
	"reflect"
)

// TODO: Still seems to deliver some false results...

func heightenIndex(cutoff int, expr interface{}) interface{} {
	switch expr := expr.(type) {
	case int:
		if expr >= cutoff {
			return expr + 1
		}

		return expr

	case LamFunc:
		res := LamFunc{}

		for _, term := range expr {
			res = append(res, lowerIndex(cutoff+1, term))
		}

		return res

	default:
		res := LamExpr{}

		for _, term := range expr.(LamExpr) {
			res = append(res, heightenIndex(cutoff, term))
		}

		return res
	}
}

func lowerIndex(cutoff int, expr interface{}) interface{} {
	switch expr := expr.(type) {
	case int:
		if expr >= cutoff {
			return expr - 1
		}

		return expr

	case LamFunc:
		res := LamFunc{}

		for _, term := range expr {
			res = append(res, lowerIndex(cutoff+1, term))
		}

		return res

	default:
		res := LamExpr{}

		for _, term := range expr.(LamExpr) {
			res = append(res, lowerIndex(cutoff, term))
		}

		return res
	}
}

// Substitute replaces index by sub
func (lx LamExpr) substitute(index int, sub interface{}) LamExpr {
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
			nw = append(nw, term.substitute(index+1, heightenIndex(0, sub)))
		}
	}

	return nw
}

// Substitute replaces index by sub
func (lf LamFunc) substitute(index int, sub interface{}) LamFunc {
	return LamFunc(LamExpr(lf).substitute(index, sub))
}

// Insert replaces index 0 by sub and returns a LamExpr
func (lf LamFunc) insert(sub interface{}) LamExpr {
	return lowerIndex(1, LamExpr(lf).substitute(0, heightenIndex(0, sub))).(LamExpr)
}

// ExpandOnce expands a lambda expression once
func (lx LamExpr) expandOnce() LamTerm {
	nw := LamExpr{}

	if len(lx) >= 2 && reflect.TypeOf(lx[0]).String() == "LamCalc.LamFunc" {
		nw = append(nw, lx[0].(LamFunc).insert(lx[1]))

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
			nw = append(nw, term.expandOnce())
		}
	}

	return nw
}

func (lf LamFunc) expandOnce() LamTerm {
	return LamFunc{LamExpr(lf).expandOnce()}
}

// Expand expands a lambda expression
func (lx LamExpr) Expand() LamFunc {
	ls := lx.simplify()
	nw := ls.expandOnce().simplify()

	for !nw.Equals(ls) {
		ls = nw
		nw = nw.expandOnce().simplify()
	}

	return nw.(LamFunc)
}

// Expand expands a lambda function (which means doing, mostly, nothing)
func (lf LamFunc) Expand() LamFunc {
	ls := lf.simplify()
	nw := ls.expandOnce().simplify()

	for !nw.Equals(ls) {
		ls = nw
		nw = nw.expandOnce().simplify()
	}

	return nw.(LamFunc)
}

// Simplify (tries) to remove unnecessary brackets
func (lx LamExpr) simplify() LamTerm {
	if len(lx) == 1 {
		if reflect.TypeOf(lx[0]).Kind() != reflect.Int {
			return lx[0].(LamTerm).simplify()
		}

		return lx

	} else if reflect.TypeOf(lx[0]).String() == "LamCalc.LamExpr" {
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
			simpl := term.simplify().(LamExpr)

			if len(simpl) == 1 {
				res = append(res, simpl[0])
			} else {
				res = append(res, simpl)
			}

		case LamFunc:
			res = append(res, term.simplify().(LamFunc))

		default:
			res = append(res, term)
		}
	}

	return res
}

// Simplify (tries) to remove unnecessary brackets
func (lf LamFunc) simplify() LamTerm {
	simpl := LamExpr(lf).simplify()

	switch simpl := simpl.(type) {
	case LamExpr:
		return LamFunc(simpl)

	default:
		return LamFunc{simpl.(LamFunc)}
	}
}
