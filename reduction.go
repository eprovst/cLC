package LamCalc

// MaxReductions determines the maximum amount of expansions before we give up
// use a negative value to have no limit (use with care...)
var MaxReductions = 10000

// Reduce reduces a lambda expression using normal order
func (lx Appl) Reduce() (Term, error) {
	return lx.NorReduce()
}

// Reduce reduces a lambda abstraction using normal order
func (la Abst) Reduce() (Term, error) {
	return la.NorReduce()
}

// Reduce returns the variable itself
func (lv Var) Reduce() (Term, error) {
	return lv.NorReduce()
}
