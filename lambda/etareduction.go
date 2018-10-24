package lambda

func (lx Appl) containsVar(idx Var) bool {
	return lx[0].containsVar(idx) || lx[1].containsVar(idx)
}

func (la Abst) containsVar(idx Var) bool {
	return la[0].containsVar(idx + 1)
}

func (lv Var) containsVar(idx Var) bool {
	return lv == idx
}

func (lf Free) containsVar(idx Var) bool {
	return false
}

// EtaReduce applies eta-reduction to the Abst when possible
func (la Abst) EtaReduce() Term {
	// First eta reduce the body
	la[0] = la[0].EtaReduce()

	switch body := la[0].(type) {
	case Appl:
		if body[1] == Var(0) && !body[0].containsVar(0) {
			// Index zero was not used anywhere else: do eta reduction.
			return lowerIndex(body[0].EtaReduce())
		}
	}

	// Else we can't do etareduction
	return la
}

// EtaReduce applies eta-reduction to the Appl when possible
func (lx Appl) EtaReduce() Term {
	lx[0] = lx[0].EtaReduce()
	lx[1] = lx[1].EtaReduce()

	return lx
}

// EtaReduce returns the variable
func (lv Var) EtaReduce() Term {
	return lv
}

// EtaReduce returns the variable
func (lf Free) EtaReduce() Term {
	return lf
}
