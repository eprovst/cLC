package lamcalc

func (lx Appl) containsVar(idx Var) bool {
	return lx[0].containsVar(idx) || lx[1].containsVar(idx)
}

func (la Abst) containsVar(idx Var) bool {
	return la[0].containsVar(idx + 1)
}

func (lv Var) containsVar(idx Var) bool {
	return lv == idx
}

// EtaReduce applies eta-reduction to the Abst when possible
func (la Abst) EtaReduce() Term {
	switch body := la[0].(type) {
	case Appl:
		if body[1] == Var(0) && !body[0].containsVar(0) {
			// Index zero was not used anywhere else: do eta reduction
			return lowerIndex(body[0].EtaReduce())
		}
	}

	// Else we can't do etareduction
	return Abst{la[0].EtaReduce()}
}

// EtaReduce applies eta-reduction to the Appl when possible
func (lx Appl) EtaReduce() Term {
	nw := Appl{}

	nw[0] = lx[0].EtaReduce()
	nw[1] = lx[1].EtaReduce()

	return nw
}

// EtaReduce returns the variable
func (lv Var) EtaReduce() Term {
	return lv
}
