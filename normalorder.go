package LamCalc

import (
	"errors"
)

// NorReduce reduces a lambda expression using normal order
func (lx LamExpr) NorReduce() (LamTerm, error) {
	ls := lx.Simplify()
	nw := ls.norReduceOnce().Simplify()

	for c := 0; !nw.alphaEquivalent(ls); c++ {
		if c == MaxReductions {
			return LamAbst{}, errors.New("exeeded maximum amount of reductions")
		}

		ls = nw
		nw = nw.norReduceOnce().Simplify()
	}

	return nw.etaReduce().Simplify(), nil
}

// norReduceOnce reduces a lambda expression once
func (lx LamExpr) norReduceOnce() LamTerm {
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
		nw = append(nw, term.norReduceOnce())
	}

	return nw
}

// NorReduce reduces a lambda abstraction using normal order
func (la LamAbst) NorReduce() (LamTerm, error) {
	ls := la.Simplify()
	nw := ls.norReduceOnce().Simplify()

	for c := 0; !nw.alphaEquivalent(ls); c++ {
		if c == MaxReductions {
			return LamAbst{}, errors.New("exeeded maximum amount of reductions")
		}

		ls = nw
		nw = nw.norReduceOnce().Simplify()
	}

	return nw.etaReduce().Simplify(), nil
}

// norReduceOnce reduces a lambda abstraction once
func (la LamAbst) norReduceOnce() LamTerm {
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
		nw = append(nw, term.norReduceOnce())
	}

	return nw
}

// NorReduce returns the variable itself
func (lv LamVar) NorReduce() (LamTerm, error) {
	return lv, nil
}

// norReduceOnce reduces a lambda variable once
func (lv LamVar) norReduceOnce() LamTerm {
	return lv
}
