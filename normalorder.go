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
			lx[1] = lx[1].norReduceOnce()
			return lx
		}

		lx[0] = lx[0].norReduceOnce()
		return lx
	}
}

// norReduceOnce reduces a lambda abstraction once
func (la Abst) norReduceOnce() Term {
	la[0] = la[0].norReduceOnce()
	return la
}

// norReduceOnce reduces a lambda variable once
func (lv Var) norReduceOnce() Term {
	return lv
}

// norReduceOnce reduces a free variable once
func (lf Free) norReduceOnce() Term {
	return lf
}

// norReduce reduces a lambda expression using normal order
func norReduce(term Term) (Term, error) {
	term = term.Copy()

	for c := 0; term.canReduce(); c++ {
		if c == MaxReductions {
			return nil, errors.New("exeeded maximum amount of reductions")
		}

		term = term.norReduceOnce()
	}

	return term.EtaReduce(), nil
}

// ConcNorReduce reduces a lambda expression using normal order,
// ignores MaxReductions, instead stops calculations once a signal is sent to done.
// If stopped mids computation puts nil on out channel.
func ConcNorReduce(term Term, out chan Term, done chan bool) {
	term = term.Copy()

	for c := 0; term.canReduce(); c++ {
		term = term.norReduceOnce()

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

// NorReduce returns the variable itself
func (lf Free) NorReduce() (Term, error) {
	return lf, nil
}
