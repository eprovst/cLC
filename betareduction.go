package LamCalc

// shiftIndex is used to correct the De Bruijn indexes
func shiftIndex(correction int, cutoff int, expr interface{}) interface{} {
	switch expr := expr.(type) {
	case int:
		if expr >= cutoff {
			return expr + correction
		}

		return expr

	case LamFunc:
		res := LamFunc{}

		for _, term := range expr {
			res = append(res, shiftIndex(correction, cutoff+1, term))
		}

		return res

	default:
		res := LamExpr{}

		for _, term := range expr.(LamExpr) {
			res = append(res, shiftIndex(correction, cutoff, term))
		}

		return res
	}
}

// substitute replaces index by sub
func (lx LamExpr) substitute(index int, sub interface{}) LamExpr {
	nw := LamExpr{}

	for _, term := range lx {
		switch term := term.(type) {
		case int:
			if term == index {
				nw = append(nw, sub)

			} else {
				nw = append(nw, term)
			}

		case LamExpr:
			nw = append(nw, term.substitute(index, sub))

		case LamFunc:
			nw = append(nw, term.substitute(index+1, shiftIndex(1, 0, sub)))
		}
	}

	return nw
}

// substitute replaces index by sub
func (lf LamFunc) substitute(index int, sub interface{}) LamFunc {
	return LamFunc(LamExpr(lf).substitute(index, sub))
}

// betaReduce replaces index 0 by sub and returns a LamExpr
func (lf LamFunc) betaReduce(sub interface{}) LamExpr {
	return shiftIndex(-1, 1, LamExpr(lf).substitute(0, shiftIndex(1, 0, sub))).(LamExpr)
}
