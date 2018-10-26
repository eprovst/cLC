package lambda

import (
	"strings"
)

// Term is a general type to represent both Applns, Absts, Vars and Frees
type Term interface {
	AlphaEquivalent(Term) bool
	EtaReduce() Term

	canReduce() bool
	containsVar(Var) bool
	substitute(Var, Term) Term

	String() string
	deDeBruijn(*strings.Builder, *[]string, *int, *[]Free)

	Reduce() (Term, error)
	NorReduce() (Term, error)
	AorReduce() (Term, error)

	norReduceOnce() Term
	aorReduceOnce() Term

	WHNF() Abst

	Copy() Term

	Serialize() string
	serialize(*strings.Builder)
}

// Appl represents an application
type Appl [2]Term

// Abst represents a lambda abstraction
type Abst [1]Term

// Var is the De Bruijn index of a bound variable minus one
type Var uint

// Free is a free variable
type Free string
