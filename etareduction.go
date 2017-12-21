package LamCalc

func (lx LamExpr) containsVar(idx LamVar) bool {
	for _, term := range lx {
		switch term := term.(type) {
		case LamVar:
			if term == idx {
				return true
			}

		case LamExpr:
			if term.containsVar(idx) {
				return true
			}

		case LamAbst:
			if term.containsVar(idx + 1) {
				return true
			}
		}
	}

	return false
}

func (la LamAbst) containsVar(idx LamVar) bool {
	return LamExpr(la).containsVar(idx)
}

func (la LamAbst) etaReduce() LamTerm {
	last := la[len(la)-1]

	if len(la) >= 2 && last == LamVar(0) {
		if !LamExpr(la[:len(la)-1]).containsVar(0) {
			// Index zero was not used anywhere else: do eta reduction
			return shiftIndex(-1, 1, LamExpr(la[:len(la)-1])).(LamExpr).etaReduce()
		}
	}

	// If zero exists it had significance: no eta reduction at this level
	return LamAbst{LamExpr(la).etaReduce()}
}

func (lx LamExpr) etaReduce() LamTerm {
	nw := LamExpr{}

	for _, term := range lx {
		switch term := term.(type) {
		case LamVar:
			nw = append(nw, term)

		case LamAbst:
			nw = append(nw, term.etaReduce())

		case LamExpr:
			nw = append(nw, term.etaReduce())
		}
	}

	return nw
}

func (lv LamVar) etaReduce() LamTerm {
	return lv
}
