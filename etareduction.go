package LamCalc

func (lx Appl) containsVar(idx Var) bool {
	return lx[0].containsVar(idx) || lx[1].containsVar(idx)
}

func (la Abst) containsVar(idx Var) bool {
	return la[0].containsVar(idx + 1)
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

	nw[0] = lx[0].etaReduce()
	nw[1] = lx[1].etaReduce()

	return nw
}

func (lv Var) etaReduce() Term {
	return lv
}
