package LamCalc

import (
	"errors"
)

// NorReduce reduces a lambda expression using normal order
func (lx Appl) NorReduce() (Term, error) {
	var ls Term = lx
	nw := ls.norReduceOnce()

	for c := 0; !nw.alphaEquivalent(ls); c++ {
		if c == MaxReductions {
			return nil, errors.New("exeeded maximum amount of reductions")
		}

		ls = nw
		nw = nw.norReduceOnce()
	}

	return nw.etaReduce(), nil
}

// norReduceOnce reduces a lambda application once
func (lx Appl) norReduceOnce() Term {
	switch fst := lx[0].(type) {
	case Abst:
		return fst.betaReduce(lx[1])

	case Var:
		return Appl{lx[0], lx[1].norReduceOnce()}

	default:
		return Appl{lx[0].norReduceOnce(), lx[1]}
	}

}

// NorReduce reduces a lambda abstraction using normal order
func (la Abst) NorReduce() (Term, error) {
	var ls Term = la
	nw := ls.norReduceOnce()

	for c := 0; !nw.alphaEquivalent(ls); c++ {
		if c == MaxReductions {
			return nil, errors.New("exeeded maximum amount of reductions")
		}

		ls = nw
		nw = nw.norReduceOnce()
	}

	return nw.etaReduce(), nil
}

// norReduceOnce reduces a lambda abstraction once
func (la Abst) norReduceOnce() Term {
	return Abst{la[0].norReduceOnce()}
}

// NorReduce returns the variable itself
func (lv Var) NorReduce() (Term, error) {
	return lv, nil
}

// norReduceOnce reduces a lambda variable once
func (lv Var) norReduceOnce() Term {
	return lv
}
