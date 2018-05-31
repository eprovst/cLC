package lamcalc

// WHNF encapsulates the expression in a lambda abstraction
func (lx Appl) WHNF() Abst {
	return Abst{Appl{heightenIndex(lx.Copy()), Var(0)}}
}

// WHNF returns the abstraction
func (la Abst) WHNF() Abst {
	// The simplification of a Abst is always a Abst
	return la.Copy().(Abst)
}

// WHNF encapsulates the Lambda variable inside of a lambda abstraction
func (lv Var) WHNF() Abst {
	return Abst{Appl{heightenIndex(lv), Var(0)}}
}
