package LamCalc

import "reflect"

// Substitute replaces index by sub
func (lf LamFunc) Substitute(index int, sub LamTerm) LamTerm {
	return LamFunc(LamExpr(lf).Substitute(index, sub).(LamExpr))
}

// Insert places sub into the function whilst consuming the latter
func (lf LamFunc) Insert(sub LamTerm) LamTerm {
	return LamExpr(lf).Substitute(0, sub)
}

// ExpandOnce expands a lambda function (which means doing nothing)
func (lf LamFunc) ExpandOnce() LamTerm {
	return lf
}

// Expand expands a lambda function (which means doing nothing)
func (lf LamFunc) Expand() LamFunc {
	return lf
}

// Substitute fills replaces index by sub
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

	return nw
}

// ExpandOnce expands a lambda expression once
func (lx LamExpr) ExpandOnce() LamTerm {
	nw := LamExpr{}

	if reflect.TypeOf(lx[0]).String() == "LamCalc.LamFunc" {
		if len(lx) < 2 {
			return lx[0].(LamFunc)

		} else if reflect.TypeOf(lx[0]).String() == "LamCalc.LamFunc" || reflect.TypeOf(lx[0]).String() == "LamCalc.LamFunc" {
			nw = append(nw, lx[0].(LamFunc).Insert(lx[1].(LamTerm)))

			if len(lx) > 2 {
				nw = append(nw, lx[2:]...)
			}

			return nw.ExpandOnce()
		}

	} else {
		for _, term := range lx {
			switch term := term.(type) {
			case int:
				nw = append(nw, term)

			case LamTerm:
				nw = append(nw, term.ExpandOnce())
			}
		}
	}

	return nw
}

// Expand expands a lambda expression
func (lx LamExpr) Expand() LamFunc {
	nw := lx.ExpandOnce()

	for reflect.TypeOf(nw).String() != "LamCalc.LamFunc" {
		nw = nw.ExpandOnce()
	}

	return nw.(LamFunc)
}
