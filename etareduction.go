package LamCalc

func (lx Appl) containsVar(idx Var) bool {
	for i := range lx {
		if lx[i].containsVar(idx) {
			return true
		}
	}

	return false
}

func (la Abst) containsVar(idx Var) bool {
	if la[0].containsVar(idx + 1) {
		return true
	}

	return false
}

func (lv Var) containsVar(idx Var) bool {
	return lv == idx
}

func (la Abst) etaReduce() Term {
	switch body := la[0].(type) {
	case Appl:
		if body[1] == Var(0) && !body[0].containsVar(0) {
			// Index zero was not used anywhere else: do eta reduction
			return shiftIndex(-1, 1, body[0].etaReduce())
		}
	}

	// Else we can't do etareduction
	return Abst{la[0].etaReduce()}
}

func (lx Appl) etaReduce() Term {
	nw := Appl{}

	for i, term := range lx {
		nw[i] = term.etaReduce()
	}

	return nw
}

func (lv Var) etaReduce() Term {
	return lv
}
