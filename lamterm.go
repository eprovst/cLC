package LamCalc

// LamTerm is a general type to represent both LamExprns and LamAbsts
type LamTerm interface {
	alphaEquivalent(LamTerm) bool

	substitute(LamVar, LamTerm) LamTerm

	containsVar(LamVar) bool
	etaReduce() LamTerm

	Equivalent(LamTerm) bool

	String() string
	deDeBruijn(boundLetters *[]string, nextletter *int) string

	Simplify() LamTerm

	Reduce() (LamTerm, error)

	NorReduce() (LamTerm, error)
	norReduceOnce() LamTerm

	WHNF() LamAbst
}

// LamExpr is a list of LamAbstns, LamExprns and De Bruijn indexes (all lowered by one) which isn't an abstraction itself.
type LamExpr []LamTerm

// LamAbst is a list of LamAbstns, LamExprns and De Bruijn indexes (all lowered by one)
type LamAbst []LamTerm

// LamVar is a variable
type LamVar uint
