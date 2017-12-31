package LamCalc

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
	case Var:
		if int(expr) >= cutoff {
			return expr + Var(correction)
		}

		return expr

	case Abst:
		res := Abst{}

		res[0] = shiftIndex(correction, cutoff+1, expr[0])

		return res

	default:
		res := Appl{}

		res[0] = shiftIndex(correction, cutoff, expr.(Appl)[0])
		res[1] = shiftIndex(correction, cutoff, expr.(Appl)[1])

		return res
	}
}
