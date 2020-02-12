// Methods for Token(s)
package main

import (
	"fmt"
)

// Nicely(?) print a []Token
func (t Tokens) String() (s string) {
	s += "["
	
	for i, v := range []Token(t) {
		s += v.String()
		
		if i < len(t) - 1 {
			s += ", "
		}
	}
	
	s += "]"
	
	return
}

// Nicely print a Token
func (t Token) String() (s string) {
	s += "{"
	
	s += t.Type.String() + ", "
	
	s += fmt.Sprint(t.Value)
	
	s += "}"

	return
}

// Give a name to a Type
func (t Type) String() (s string) {
	switch t {
	case Number:		s = "Number"
	case OpenParen:		s = "OpenParen"
	case CloseParen:	s = "CloseParen"
	case Multiply:		s = "Multiply"
	case Divide:		s = "Divide"
	case Add:			s = "Add"
	case Subtract:		s = "Subtract"
	default:			s = "<nil>"
	}
	
	return
}
