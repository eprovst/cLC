package lamcalc

import (
	"errors"
)

// norReduceOnce reduces a lambda application once
func (lx Appl) norReduceOnce() Term {
	switch fst := lx[0].(type) {
	case Abst:
		return fst.BetaReduce(lx[1])

	default:
		if !lx[0].canReduce() {
			return Appl{lx[0], lx[1].norReduceOnce()}
		}

		return Appl{lx[0].norReduceOnce(), lx[1]}
	}

}

// norReduceOnce reduces a lambda abstraction once
func (la Abst) norReduceOnce() Term {
	return Abst{la[0].norReduceOnce()}
}

// norReduceOnce reduces a lambda variable once
func (lv Var) norReduceOnce() Term {
	return lv
}

// norReduce reduces a lambda expression using normal order
func norReduce(term Term) (Term, error) {
	for c := 0; term.canReduce(); c++ {
		if c == MaxReductions {
			return nil, errors.New("exeeded maximum amount of reductions")
		}

		term = term.norReduceOnce()
	}

	return term.EtaReduce(), nil
}

// NorReduce reduces an application using normal order
func (lx Appl) NorReduce() (Term, error) {
	return norReduce(lx)
}

// NorReduce reduces a lambda abstraction using normal order
func (la Abst) NorReduce() (Term, error) {
	return norReduce(la)
}

// NorReduce returns the variable itself
func (lv Var) NorReduce() (Term, error) {
	return lv, nil
}
