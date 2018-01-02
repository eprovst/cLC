package LamCalc

import (
	"errors"
)

// aorReduceOnce reduces a lambda application once
func (lx Appl) aorReduceOnce() Term {
	if !lx[1].canReduce() {
		switch fst := lx[0].(type) {
		case Abst:
			return fst.BetaReduce(lx[1])

		default:
			return Appl{lx[0].aorReduceOnce(), lx[1]}
		}
	}

	return Appl{lx[0], lx[1].aorReduceOnce()}
}

// aorReduceOnce reduces a lambda abstraction once
func (la Abst) aorReduceOnce() Term {
	return Abst{la[0].aorReduceOnce()}
}

// aorReduceOnce reduces a lambda variable once
func (lv Var) aorReduceOnce() Term {
	return lv
}

// aorReduce reduces a lambda expression using applicative order
func aorReduce(term Term) (Term, error) {
	for c := 0; term.canReduce(); c++ {
		if c == MaxReductions {
			return nil, errors.New("exeeded maximum amount of reductions")
		}

		term = term.aorReduceOnce()
	}

	return term.EtaReduce(), nil
}

// AorReduce reduces an application using applicative order
func (lx Appl) AorReduce() (Term, error) {
	return aorReduce(lx)
}

// AorReduce reduces a lambda abstraction using applicative order
func (la Abst) AorReduce() (Term, error) {
	return aorReduce(la)
}

// AorReduce returns the variable itself
func (lv Var) AorReduce() (Term, error) {
	return lv, nil
}
