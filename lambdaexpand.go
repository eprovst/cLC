package LamCalc

import "reflect"

// TODO: Do a lot of testing to see if this acctualy works...

// Substitute replaces index by sub
func (lf LamFunc) Substitute(index int, sub LamTerm) LamTerm {
	nw := LamFunc{}
	for _, term := range lf {
		switch term := term.(type) {
		case int:
			if term == index {
				nw = append(nw, sub)
			} else {
				nw = append(nw, term)
			}

		case LamExpr:
			nw = append(nw, term.Substitute(index, sub))

		case LamFunc:
			nw = append(nw, term.Substitute(index+1, sub))

		}
	}

	return nw
}

// Insert places sub into the function whilst consuming the latter
func (lf LamFunc) Insert(sub LamTerm) LamTerm {
	nw := LamExpr{}
	for _, term := range lf {
		switch term := term.(type) {
		case int:
			if term == 0 {
				nw = append(nw, sub)
			} else {
				nw = append(nw, term)
			}

		case LamExpr:
			nw = append(nw, term.Substitute(0, sub))

		case LamFunc:
			nw = append(nw, term.Substitute(1, sub))

		}
	}

	return nw
}

// ExpandOnce expands a lambda function (which means doing nothing)
func (lf LamFunc) ExpandOnce() LamTerm {
	return lf
}

// Expand expands a lambda function (which means doing nothing)
func (lf LamFunc) Expand() LamFunc {
	return lf
}

// Substitute fills relpaces index by sub
func (lx LamExpr) Substitute(index int, sub LamTerm) LamTerm {
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
			nw = append(nw, term.Substitute(index, sub))

		case LamFunc:
			nw = append(nw, term.Substitute(index+1, sub))

		}
	}

	if len(nw) == 1 {
		return nw[0].(LamTerm)
	}

	return nw
}

// ExpandOnce expands a lambda expression once
func (lx LamExpr) ExpandOnce() LamTerm {
	nw := LamExpr{}

	switch term := lx[0].(type) {
	case LamExpr:
		nw = append(nw, term.Expand())
		if len(lx) > 1 {
			nw = append(nw, lx[1:])
		}

		return nw

	case LamFunc:
		if len(lx) > 1 {
			nw = append(nw, term.Insert(lx[1].(LamTerm)))
			if len(lx) > 2 {
				return append(nw, lx[2:])
			}

			return nw[0].(LamExpr)
		}

		return term

	default:
		return lx
	}
}

// Expand expands a lambda expression
func (lx LamExpr) Expand() LamFunc {
	nw := lx.ExpandOnce()

	for reflect.TypeOf(nw).String() != "LamCalc.LamFunc" {
		nw = nw.ExpandOnce()
	}

	return nw.(LamFunc)
}
