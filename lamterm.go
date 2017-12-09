package LamCalc

// LamTerm is a general type to represent both LamExprns and LamFuncs
type LamTerm interface {
	// Allow us to manipulate it as a list
	len() int
	index(int) interface{}
	append(...interface{}) LamTerm

	// equivalence.go
	Equivalent(LamTerm) bool

	// alphaequivalence.go
	alphaEquivalent(LamTerm) bool

	// stringify.go
	String() string
	deDeBruijn(boundLetters []string, nextletter int) string

	// simplify.go
	simplify() LamTerm

	// reduction.go
	Reduce() LamFunc
	reduceOnce() LamTerm
}

// LamExpr is a list of lamfuncs, lamexprns and De Bruijn indexes (all lowered by one) which isn't a function itself.
type LamExpr []interface{}

// Utils for LamExpr
func (lx LamExpr) len() int {
	return len(lx)
}

func (lx LamExpr) index(i int) interface{} {
	return lx[i]
}

func (lx LamExpr) append(terms ...interface{}) LamTerm {
	return append(lx, terms...)
}

// LamFunc is a list of lamfuncs, lamexprns and De Bruijn indexes (all lowered by one)
type LamFunc []interface{}

// Utils for LamFunc
func (lf LamFunc) len() int {
	return len(lf)
}

func (lf LamFunc) index(i int) interface{} {
	return lf[i]
}

func (lf LamFunc) append(terms ...interface{}) LamTerm {
	return append(lf, terms...)
}
