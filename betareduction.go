package lamcalc

// substitute replaces index by sub
func (lx Appl) substitute(index Var, sub Term) Term {
	lx[0] = lx[0].substitute(index, sub)
	lx[1] = lx[1].substitute(index, sub)

	return lx
}

// substitute replaces index by sub
func (la Abst) substitute(index Var, sub Term) Term {
	la[0] = la[0].substitute(index+1, heightenIndex(sub))
	return la
}

// substitute replaces index by sub
func (lv Var) substitute(index Var, sub Term) Term {
	if lv == index {
		return sub.Copy()
	}

	return lv
}

// substitute replaces index by sub
func (lf Free) substitute(index Var, sub Term) Term {
	return lf
}

// BetaReduce substitutes the localy bound variable of the Abst by sub
func (la Abst) BetaReduce(sub Term) Term {
	return lowerIndex(
		la[0].substitute(0, heightenIndex(sub)),
	)
}
