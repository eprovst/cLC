package LamCalc

import "bytes"

// Term is a general type to represent both Applns and Absts
type Term interface {
	alphaEquivalent(Term) bool

	substitute(Var, Term) Term

	containsVar(Var) bool
	etaReduce() Term

	Equivalent(Term) bool

	String() string
	deDeBruijn(*bytes.Buffer, *[]string, *int)

	canReduce() bool
	Reduce() (Term, error)

	NorReduce() (Term, error)
	norReduceOnce() Term

	AorReduce() (Term, error)
	aorReduceOnce() Term

	WHNF() Abst
}

// Appl is a list of Abstns, Applns and De Bruijn indexes (all lowered by one) which isn't an abstraction itself.
type Appl [2]Term

// Abst is a list of Abstns, Applns and De Bruijn indexes (all lowered by one)
type Abst [1]Term

// Var is a variable
type Var uint
