package LamCalc

func (lx LamExpr) containsVar(idx LamVar) bool {
	for _, term := range lx {
		switch term := term.(type) {
		case LamAbst:
			if term.containsVar(idx + 1) {
				return true
			}

		default:
			if term.containsVar(idx) {
				return true
			}
		}
	}

	return false
}

func (la LamAbst) containsVar(idx LamVar) bool {
	for _, term := range la {
		switch term := term.(type) {
		case LamAbst:
			if term.containsVar(idx + 1) {
				return true
			}

		default:
			if term.containsVar(idx) {
				return true
			}
		}
	}

	return false
}

func (lv LamVar) containsVar(idx LamVar) bool {
	return lv == idx
}

func (la LamAbst) etaReduce() LamTerm {
	last := la[len(la)-1]

	if len(la) >= 2 && last == LamVar(0) {
		if !LamExpr(la[:len(la)-1]).containsVar(0) {
			// Index zero was not used anywhere else: do eta reduction
			return shiftIndex(-1, 1, LamExpr(la[:len(la)-1])).etaReduce()
		}
	}

	// If zero exists it had significance: no eta reduction at this level
	nw := LamAbst{}

	for _, term := range la {
		nw = append(nw, term.etaReduce())
	}

	return nw
}

func (lx LamExpr) etaReduce() LamTerm {
	nw := LamExpr{}

	for _, term := range lx {
		nw = append(nw, term.etaReduce())
	}

	return nw
}

func (lv LamVar) etaReduce() LamTerm {
	return lv
}
