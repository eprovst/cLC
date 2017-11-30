package LamCalc

import "reflect"

// LamTerm is a general type to represent both LamExprns and LamFuncs.
type LamTerm interface {
	Equals(other LamTerm) bool
	String() string
	deDebruijn(boundLetters []string, nextletter int) string

	// Allow us to manipulate it as a list
	Len() int
	Index(int) interface{}
}

// LamExpr is a list of lamfuncs, lamexprns and De Bruijn indexes (all lowered by one) which isn't a function itself.
type LamExpr []interface{}

// Len returns the lenght of the LamExpr
func (lx LamExpr) Len() int {
	return len(lx)
}

// Index returns the ith element of the LamExpr
func (lx LamExpr) Index(i int) interface{} {
	return lx[i]
}

// Equals checks wether the LamExp is identical to a LamTerm
func (lx LamExpr) Equals(other LamTerm) bool {

	if reflect.TypeOf(other).String() != "LamCalc.LamExpr" {
		return false

	} else if len(lx) != other.Len() {
		return false
	}

	for i := range lx {
		switch elem := lx[i].(type) {
		case int:
			if reflect.TypeOf(other.Index(i)).Kind() != reflect.Int || elem != other.Index(i).(int) {
				return false
			}

		case LamExpr:
			if reflect.TypeOf(other.Index(i)).String() != "LamCalc.LamExpr" || !elem.Equals(other.Index(i).(LamExpr)) {
				return false
			}

		case LamFunc:
			if reflect.TypeOf(other.Index(i)).String() != "LamCalc.LamFunc" || !LamExpr(elem).Equals(LamExpr(other.Index(i).(LamFunc))) {
				return false
			}

		default:
			return false
		}
	}

	return true
}

// LamFunc is a list of lamfuncs, lamexprns and De Bruijn indexes (all lowered by one)
type LamFunc []interface{}

// Len returns the lenght of the LamFunc
func (lf LamFunc) Len() int {
	return len(lf)
}

// Index returns the ith element of the LamFunc
func (lf LamFunc) Index(i int) interface{} {
	return lf[i]
}

// Equals checks wether two LamFuncs are identical
func (lf LamFunc) Equals(other LamTerm) bool {

	if reflect.TypeOf(other).String() != "LamCalc.LamFunc" {
		return false
	}

	return LamExpr(lf).Equals(LamExpr(other.(LamFunc)))
}
