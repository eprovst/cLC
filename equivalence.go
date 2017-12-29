package LamCalc

// Currently someterm.Equals(otherterm) only uses alphaEquivalence
// should this include other checks?

// Equivalent checks if the Appl and Term are Equivalent
func (lx Appl) Equivalent(other Term) bool {
	return lx.alphaEquivalent(other)
}

// Equivalent checks if the Abst and Term are Equivalent
func (la Abst) Equivalent(other Term) bool {
	return la.alphaEquivalent(other)
}

// Equivalent checks if the Var and Term are Equivalent
func (lv Var) Equivalent(other Term) bool {
	return lv.alphaEquivalent(other)
}
