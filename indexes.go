package LamCalc

// heightenIndex heightens the De Bruijn indexes by one where needed
func heightenIndex(expr LamTerm) LamTerm {
	return shiftIndex(1, 0, expr)
}

// lowerIndex lowers the De Bruijn indexes by one where needed
func lowerIndex(expr LamTerm) LamTerm {
	return shiftIndex(-1, 1, expr)
}

// shiftIndex is used to correct the De Bruijn indexes
func shiftIndex(correction int, cutoff int, expr LamTerm) LamTerm {
	switch expr := expr.(type) {
	case LamVar:
		if int(expr) >= cutoff {
			return expr + LamVar(correction)
		}

		return expr

	case LamAbst:
		res := LamAbst{}

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
