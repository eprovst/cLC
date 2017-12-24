package LamCalc

import (
	"errors"
)

// MaxReductions determines the maximum amount of expansions before we give up
// use a negative value to have no limit (use with care...)
var MaxReductions = 10000

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

// reduceOnce reduces a lambda expression once
func (lx LamExpr) reduceOnce() LamTerm {
	nw := LamExpr{}

	switch fst := lx[0].(type) {
	case LamAbst:
		if len(lx) >= 2 {
			nw = append(nw, fst.betaReduce(lx[1]))

			if len(lx) > 2 {
				nw = append(nw, lx[2:]...)
			}

			return nw
		}
	}

	for _, term := range lx {
		nw = append(nw, term.reduceOnce())
	}

	return nw
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

// reduceOnce reduces a lambda abstraction once
func (la LamAbst) reduceOnce() LamTerm {
	nw := LamAbst{}

	switch fst := la[0].(type) {
	case LamAbst:
		if len(la) >= 2 {
			nw = append(nw, fst.betaReduce(la[1]))

			if len(la) > 2 {
				nw = append(nw, la[2:]...)
			}

			return nw
		}
	}

	for _, term := range la {
		nw = append(nw, term.reduceOnce())
	}

	return nw
}

// Reduce returns the variable itself
func (lv LamVar) Reduce() (LamTerm, error) {
	return lv, nil
}

// reduceOnce reduces a lambda variable once
func (lv LamVar) reduceOnce() LamTerm {
	return lv
}
