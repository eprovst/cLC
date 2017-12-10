package LamCalc

// shiftIndex is used to correct the De Bruijn indexes
func shiftIndex(correction int, cutoff int, expr interface{}) interface{} {
	switch expr := expr.(type) {
	case int:
		if expr >= cutoff {
			return expr + correction
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
