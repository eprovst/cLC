package lamcalc

import "bytes"

// Term is a general type to represent both Applns, Absts and Vars
type Term interface {
	AlphaEquivalent(Term) bool
	EtaReduce() Term

	canReduce() bool
	containsVar(Var) bool
	substitute(Var, Term) Term

	String() string
	deDeBruijn(*bytes.Buffer, *[]string, *int)

	Reduce() (Term, error)
	NorReduce() (Term, error)
	AorReduce() (Term, error)

	norReduceOnce() Term
	aorReduceOnce() Term

	WHNF() Abst

	Serialize() string
	serialize(*bytes.Buffer)
}

// Appl represents an application
type Appl [2]Term

// Abst represents a lambda abstraction
type Abst [1]Term

// Var is the De Bruijn index of a variable minus one
type Var uint
