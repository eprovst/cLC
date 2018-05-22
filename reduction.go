package lamcalc

// MaxReductions determines the maximum amount of expansions before we give up
// use a negative value to have no limit (use with care...)
var MaxReductions = 10000

// Checks if the term can be reduced (only checks for beta reduction, not eta)
func (lx Appl) canReduce() bool {
	switch lx[0].(type) {
	case Abst:
		return true

	default:
		return lx[0].canReduce() || lx[1].canReduce()
	}
}

func (la Abst) canReduce() bool {
	return la[0].canReduce()
}

func (lv Var) canReduce() bool {
	return false
}

// Reduce reduces an application
func (lx Appl) Reduce() (Term, error) {
	return aorReduce(lx)
}

// Reduce reduces a lambda abstraction
func (la Abst) Reduce() (Term, error) {
	return aorReduce(la)
}

// Reduce returns the variable itself
func (lv Var) Reduce() (Term, error) {
	return aorReduce(lv)
}
