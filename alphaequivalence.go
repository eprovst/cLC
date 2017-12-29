package LamCalc

// As we are using De Bruijn indexes alphaEquivalence is the same as syntactic equivalence

// Equals checks wether the LamExp is identical to a Term
func (lx Appl) alphaEquivalent(other Term) bool {
	switch other := other.(type) {
	case Appl:
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

// Equals checks wether a Abst and a Term are identical
func (la Abst) alphaEquivalent(other Term) bool {
	switch other := other.(type) {
	case Abst:
		if !la[0].alphaEquivalent(other[0]) {
			return false
		}

		return true

	default:
		return false
	}
}

// Equals checks wether a Var and a Term are identical
func (lv Var) alphaEquivalent(other Term) bool {
	switch other := other.(type) {
	case Var:
		return lv == other

	default:
		return false
	}
}
