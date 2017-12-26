package LamCalc

// MaxReductions determines the maximum amount of expansions before we give up
// use a negative value to have no limit (use with care...)
var MaxReductions = 10000

// Reduce reduces a lambda expression using normal order
func (lx LamExpr) Reduce() (LamTerm, error) {
	return lx.NorReduce()
}

// Reduce reduces a lambda abstraction using normal order
func (la LamAbst) Reduce() (LamTerm, error) {
	return la.NorReduce()
}

// Reduce returns the variable itself
func (lv LamVar) Reduce() (LamTerm, error) {
	return lv.NorReduce()
}
