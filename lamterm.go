package LamCalc

// TODO: Use the fact that LamVar is a LamTerm more in the source

// LamTerm is a general type to represent both LamExprns and LamAbsts
type LamTerm interface {
	// alphaequivalence.go
	alphaEquivalent(LamTerm) bool

	// betareduction.go
	substitute(LamVar, LamTerm) LamTerm

	// etareduction.go
	etaReduce() LamTerm

	// equivalence.go
	Equivalent(LamTerm) bool

	// stringify.go
	String() string
	deDeBruijn(boundLetters []string, nextletter LamVar) string

	// simplify.go
	Simplify() LamTerm

	// reduction.go
	Reduce() (LamTerm, error)
	WHNFReduce() (LamAbst, error)
	reduceOnce() LamTerm
}

// LamExpr is a list of LamAbstns, LamExprns and De Bruijn indexes (all lowered by one) which isn't an abstraction itself.
type LamExpr []LamTerm

// LamAbst is a list of LamAbstns, LamExprns and De Bruijn indexes (all lowered by one)
type LamAbst []LamTerm

// LamVar is a bound variable
type LamVar uint
