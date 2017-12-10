package LamCalc

import "reflect"

func (lx LamExpr) containsIdx(idx int) bool {
	for _, term := range lx {
		switch term := term.(type) {
		case int:
			if term == idx {
				return true
			}

		case LamExpr:
			if term.containsIdx(idx) {
				return true
			}

		case LamAbst:
			if term.containsIdx(idx + 1) {
				return true
			}
		}
	}

	return false
}

func (lf LamAbst) containsIdx(idx int) bool {
	return LamExpr(lf).containsIdx(idx)
}

func (lf LamAbst) etaReduce() LamTerm {
	last := lf[len(lf)-1]

	if len(lf) >= 2 && reflect.TypeOf(last).Kind() == reflect.Int && last.(int) == 0 {
		if !LamExpr(lf[:len(lf)-1]).containsIdx(0) {
			// Index zero was not used anywhere else: do eta reduction
			return shiftIndex(-1, 1, LamExpr(lf[:len(lf)-1])).(LamExpr).etaReduce()
		}
	}

	// If zero exists it had significance: no eta reduction at this level
	return LamAbst{LamExpr(lf).etaReduce()}
}

func (lx LamExpr) etaReduce() LamTerm {
	nw := LamExpr{}

	for _, term := range lx {
		switch term := term.(type) {
		case int:
			nw = append(nw, term)

		case LamAbst:
			nw = append(nw, term.etaReduce())

		case LamExpr:
			nw = append(nw, term.etaReduce())
		}
	}

	return nw
}
