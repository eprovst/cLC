package lamcalc

// Copy creates a copy of the application
func (lx Appl) Copy() Term {
	return Appl{lx[0].Copy(), lx[1].Copy()}
}

// Copy creates a copy of the abstraction
func (la Abst) Copy() Term {
	return Abst{la[0].Copy()}
}

// Copy returns the value of the variable
func (lv Var) Copy() Term {
	return lv
}
