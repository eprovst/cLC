package lamcalc

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
			lx[0] = lx[0].aorReduceOnce()
			return lx
		}
	}

	lx[1] = lx[1].aorReduceOnce()
	return lx
}

// aorReduceOnce reduces a lambda abstraction once
func (la Abst) aorReduceOnce() Term {
	la[0] = la[0].aorReduceOnce()
	return la
}

// aorReduceOnce reduces a lambda variable once
func (lv Var) aorReduceOnce() Term {
	return lv
}

// aorReduceOnce reduces a free variable once
func (lf Free) aorReduceOnce() Term {
	return lf
}

// aorReduce reduces a lambda expression using applicative order
func aorReduce(term Term) (Term, error) {
	term = term.Copy()

	for c := 0; term.canReduce(); c++ {
		if c == MaxReductions {
			return nil, errors.New("exeeded maximum amount of reductions")
		}

		term = term.aorReduceOnce()
	}

	return term.EtaReduce(), nil
}

// ConcAorReduce reduces a lambda expression using applicative order,
// ignores MaxReductions, instead stops calculations once a signal is sent to done.
// If stopped mids computation puts nil on out channel.
func ConcAorReduce(term Term, out chan Term, done chan bool) {
	term = term.Copy()

	for c := 0; term.canReduce(); c++ {
		term = term.aorReduceOnce()

		// Stop if a signal is sent to done
		select {
		case <-done:
			out <- nil
			return
		default:
		}
	}

	out <- term.EtaReduce()
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

// AorReduce returns the variable itself
func (lf Free) AorReduce() (Term, error) {
	return lf, nil
}
