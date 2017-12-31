package LamCalc

// substitute replaces index by sub
func (lx Appl) substitute(index Var, sub Term) Term {
	nw := Appl{}

	nw[0] = lx[0].substitute(index, sub)
	nw[1] = lx[1].substitute(index, sub)

	return nw
}

// substitute replaces index by sub
func (la Abst) substitute(index Var, sub Term) Term {
	return Abst{la[0].substitute(index+1, heightenIndex(sub))}
}

// substitute replaces index by sub
func (lv Var) substitute(index Var, sub Term) Term {
	if lv == index {
		return sub
	}

	return lv
}

// betaReduce replaces index 0 by sub and returns a Appl
func (la Abst) betaReduce(sub Term) Term {
	return lowerIndex(
		la[0].substitute(0, heightenIndex(sub)),
	)
}
