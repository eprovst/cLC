package lamcalc

// heightenIndex heightens the De Bruijn indexes by one where needed
func heightenIndex(expr Term) Term {
	return shiftIndex(1, 0, expr)
}

// lowerIndex lowers the De Bruijn indexes by one where needed
func lowerIndex(expr Term) Term {
	return shiftIndex(-1, 1, expr)
}

// shiftIndex is used to correct the De Bruijn indexes
func shiftIndex(correction int, cutoff int, expr Term) Term {
	switch expr := expr.(type) {

	case Abst:
		expr[0] = shiftIndex(correction, cutoff+1, expr[0])

		return expr

	case Appl:
		expr[0] = shiftIndex(correction, cutoff, expr[0])
		expr[1] = shiftIndex(correction, cutoff, expr[1])

		return expr

	case Var:
		if int(expr) >= cutoff {
			return expr + Var(correction)
		}

		return expr

	default:
		return expr
	}
}
