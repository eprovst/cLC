package LamCalc

// As we are using De Bruijn indexes alphaEquivalence is the same as syntactic equivalence

// Equals checks wether the LamExp is identical to a LamTerm
func (lx LamExpr) alphaEquivalent(other LamTerm) bool {
	switch other := other.(type) {
	case LamExpr:
		if len(lx) != len(other) {
			return false
		}

		for i := range lx {
			if !lx[i].alphaEquivalent(other[i]) {
				return false
			}
		}

		return true

	default:
		return false
	}
}

// Equals checks wether a LamAbst and a LamTerm are identical
func (la LamAbst) alphaEquivalent(other LamTerm) bool {
	switch other := other.(type) {
	case LamAbst:
		if len(la) != len(other) {
			return false
		}

		for i := range la {
			if !la[i].alphaEquivalent(other[i]) {
				return false
			}
		}

		return true

	default:
		return false
	}
}

// Equals checks wether a LamVar and a LamTerm are identical
func (lv LamVar) alphaEquivalent(other LamTerm) bool {
	switch other := other.(type) {
	case LamVar:
		return lv == other

	default:
		return false
	}
}
