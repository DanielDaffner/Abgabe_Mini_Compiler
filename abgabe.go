package main

import (
	"fmt"
	"unicode"
)
import "strconv"

// Interface

type Exp interface {
	pretty() string
	eval(s ValState) Val
	infer(t TyState) (Type, int)
}

// Statement

type Stmt interface {
	pretty() string
	eval(s ValState)
	check(t TyState) (bool, int, int)
}

var varName string
var inputLength int
var errorLength int

type Bool bool
type Num int
type Mult [2]Exp
type Plus [2]Exp
type And [2]Exp
type Or [2]Exp
type Neg [1]Exp
type Equ [2]Exp
type Les [2]Exp
type Var string

type Block struct {
	s Stmt
}
type ComS [2]Stmt
type Decl struct {
	lhs string
	rhs Exp
}
type Assign struct {
	name  string
	value Exp
}
type While struct {
	e Exp
	b Block
}
type IfEl struct {
	e  Exp
	b1 Block
	b2 Block
}
type Print struct {
	e Exp
}

type ValState map[string]Val
type TyState map[string]Type

// Values

type Kind int

const (
	ValueInt  Kind = 0
	ValueBool Kind = 1
	Undefined Kind = 2
)

type Val struct {
	flag Kind
	valI int
	valB bool
}

func mkInt(x int) Val {
	return Val{flag: ValueInt, valI: x}
}
func mkBool(x bool) Val {
	return Val{flag: ValueBool, valB: x}
}
func mkUndefined() Val {
	return Val{flag: Undefined}
}

func showVal(v Val) string {
	var s string
	switch {
	case v.flag == ValueInt:
		s = Num(v.valI).pretty()
	case v.flag == ValueBool:
		s = Bool(v.valB).pretty()
	case v.flag == Undefined:
		s = "Undefined"
	}
	return s
}

// Types

type Type int

const (
	TyIllTyped Type = 0
	TyInt      Type = 1
	TyBool     Type = 2
)

const (
	Integer        int = 1
	Booleans       int = 2
	Addition       int = 3
	Multiplication int = 4
	Disjunction    int = 5
	Conjuction     int = 6
	Negation       int = 7
	Equality       int = 8
	Lesser         int = 9
	Variables      int = 10
	Condition      int = 11
	BlockT         int = 12
)

func showType(t Type) string {
	var s string
	switch {
	case t == TyInt:
		s = "Int"
	case t == TyBool:
		s = "Bool"
	case t == TyIllTyped:
		s = "Illtyped"
	}
	return s
}

// pretty print

func (x Bool) pretty() string {
	if x {
		return "true"
	} else {
		return "false"
	}

}

func (x Num) pretty() string {
	return strconv.Itoa(int(x))
}

func (e Mult) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "*"
	x += e[1].pretty()
	x += ")"

	return x
}

func (e Plus) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "+"
	x += e[1].pretty()
	x += ")"

	return x
}

func (e And) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "&&"
	x += e[1].pretty()
	x += ")"

	return x
}

func (e Or) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "||"
	x += e[1].pretty()
	x += ")"

	return x
}

// Negation
func (e Neg) pretty() string {
	var x string
	x = "!"
	x += e[0].pretty()
	return x
}

// Equality
func (e Equ) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "=="
	x += e[1].pretty()
	x += ")"

	return x
}

// Lesser Test
func (e Les) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "<"
	x += e[1].pretty()
	x += ")"

	return x
}

// Vars

func (x Var) pretty() string {
	return (string)(x)
}

// Command Sequence
func (s ComS) pretty() string {
	var x string
	x = s[0].pretty()
	x += " ; "
	x += s[1].pretty()
	return x
}

// Variable declaration
func (e Decl) pretty() string {
	var x string
	x = e.lhs
	x += " := "
	x += e.rhs.pretty()
	return x
}

// Variable assignment
func (e Assign) pretty() string {
	var x string
	x = e.name
	x += " = "
	x += e.value.pretty()
	return x
}

// While
func (w While) pretty() string {
	var x string
	x = " while "
	x += w.e.pretty()
	x += " { "
	x += w.b.pretty()
	x += " } "
	return x
}

// If-then-else

func (ifel IfEl) pretty() string {
	var x string
	x = "if "
	x += ifel.e.pretty()
	x += " then "
	x += ifel.b1.pretty()
	x += " else "
	x += ifel.b2.pretty()
	return x
}

// Print

func (e Print) pretty() string {
	var x string
	x = "print: "
	x += e.e.pretty()
	return x
}

// Block

func (b Block) pretty() string {
	var x string
	x = b.s.pretty()
	return x
}

// Val

func (v Val) pretty() string {
	var x string
	switch v.flag {
	case ValueInt:
		x = strconv.Itoa(v.valI)
		return x
	case ValueBool:
		x = strconv.FormatBool(v.valB)
		return x
	default:
		x = "illtyped"
		return x
	}
}

// Evaluator

func (x Bool) eval(s ValState) Val {
	return mkBool((bool)(x))
}

func (x Num) eval(s ValState) Val {
	return mkInt((int)(x))
}

func (e Mult) eval(s ValState) Val {
	n1 := e[0].eval(s)
	n2 := e[1].eval(s)
	if n1.flag == ValueInt && n2.flag == ValueInt {
		return mkInt(n1.valI * n2.valI)
	}
	return mkUndefined()
}

func (e Plus) eval(s ValState) Val {
	n1 := e[0].eval(s)
	n2 := e[1].eval(s)
	if n1.flag == ValueInt && n2.flag == ValueInt {
		return mkInt(n1.valI + n2.valI)
	}
	return mkUndefined()
}

func (e And) eval(s ValState) Val {
	b1 := e[0].eval(s)
	b2 := e[1].eval(s)
	switch {
	case b1.flag == ValueBool && b1.valB == false:
		return mkBool(false)
	case b1.flag == ValueBool && b2.flag == ValueBool:
		return mkBool(b1.valB && b2.valB)
	}
	return mkUndefined()
}

func (e Or) eval(s ValState) Val {
	b1 := e[0].eval(s)
	b2 := e[1].eval(s)
	switch {
	case b1.flag == ValueBool && b1.valB == true:
		return mkBool(true)
	case b1.flag == ValueBool && b2.flag == ValueBool:
		return mkBool(b1.valB || b2.valB)
	}
	return mkUndefined()
}

// Negation

func (e Neg) eval(s ValState) Val {
	b1 := e[0].eval(s)
	if b1.flag == ValueBool {
		return mkBool(!b1.valB)
	}
	return mkUndefined()
}

// Equality Test

func (e Equ) eval(s ValState) Val {
	b1 := e[0].eval(s)
	b2 := e[1].eval(s)
	switch {
	case b1.flag == ValueBool && b2.flag == ValueBool:
		if b1.valB == b2.valB {
			return mkBool(true)
		}
		return mkBool(false)
	case b1.flag == ValueInt && b2.flag == ValueInt:
		if b1.valI == b2.valI {
			return mkBool(true)
		}
		return mkBool(false)
	}
	return mkUndefined()
}

// Lesser Test

func (e Les) eval(s ValState) Val {
	b1 := e[0].eval(s)
	b2 := e[1].eval(s)
	if b1.flag == ValueInt && b2.flag == ValueInt {
		if b1.valI < b2.valI {
			return mkBool(true)
		}
		return mkBool(false)
	}
	return mkUndefined()
}

// vars

func (x Var) eval(s ValState) Val {
	return s[(string)(x)]
}

// Exp

// Command Sequence
func (x ComS) eval(s ValState) {
	x[0].eval(s)
	x[1].eval(s)
}

// Variable declaration
func (decl Decl) eval(s ValState) {
	v := decl.rhs.eval(s)
	x := (string)(decl.lhs)
	s[x] = v

}

// Variable assignment
func (assign Assign) eval(s ValState) {
	v, ok := s[assign.name]
	if ok {
		v = assign.value.eval(s)
		s[assign.name] = v
	}
}

// While
func (w While) eval(s ValState) {
	if w.e.eval(s).valB {
		w.b.eval(s)
		w.eval(s)
	}
}

// If-then-else
func (ifel IfEl) eval(s ValState) {
	if ifel.e.eval(s).valB {
		ifel.b1.eval(s)
	} else {
		ifel.b2.eval(s)
	}
}

// Print
func (p Print) eval(s ValState) {
	p1 := p.e.eval(s)
	switch p1.flag {
	case ValueInt:
		fmt.Printf("\n %d", p1.valI)
	case ValueBool:
		fmt.Printf("\n %t", p1.valB)
	}

}

// Block

func (b Block) eval(s ValState) {
	b.s.eval(s)
}

// Type inferencer/checker

func (x Bool) infer(t TyState) (Type, int) {
	return TyBool, Booleans
}

func (x Num) infer(t TyState) (Type, int) {
	return TyInt, Integer
}

func (e Mult) infer(t TyState) (Type, int) {
	t1, _ := e[0].infer(t)
	t2, _ := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyInt, Multiplication
	}
	return TyIllTyped, Multiplication
}

func (e Plus) infer(t TyState) (Type, int) {
	t1, _ := e[0].infer(t)
	t2, _ := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyInt, Addition
	}
	return TyIllTyped, Addition
}

func (e And) infer(t TyState) (Type, int) {
	t1, _ := e[0].infer(t)
	t2, _ := e[1].infer(t)
	if t1 == TyBool && t2 == TyBool {
		return TyBool, Conjuction
	}
	return TyIllTyped, Conjuction
}

func (e Or) infer(t TyState) (Type, int) {
	t1, _ := e[0].infer(t)
	t2, _ := e[1].infer(t)
	if t1 == TyBool && t2 == TyBool {
		return TyBool, Disjunction
	}
	return TyIllTyped, Disjunction
}

// Negation
func (e Neg) infer(t TyState) (Type, int) {
	t1, _ := e[0].infer(t)

	if t1 == TyBool {
		return TyBool, Negation
	}
	return TyIllTyped, Negation
}

// Equality Test
func (e Equ) infer(t TyState) (Type, int) {
	t1, _ := e[0].infer(t)
	t2, _ := e[1].infer(t)
	if t1 == TyBool && t2 == TyBool {
		return TyBool, Equality
	}
	if t1 == TyInt && t2 == TyInt {
		return TyBool, Equality
	}
	return TyIllTyped, Equality
}

//Lesser Test
func (e Les) infer(t TyState) (Type, int) {
	t1, _ := e[0].infer(t)
	t2, _ := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyBool, Lesser
	}
	return TyIllTyped, Lesser
}

// Vars

func (x Var) infer(t TyState) (Type, int) {
	y := (string)(x)
	ty, ok := t[y]
	if ok {
		return ty, Variables
	} else {
		return TyIllTyped, Variables
	}

}

// Check decl

func (e Decl) check(t TyState) (bool, int, int) {
	v, vP := e.rhs.infer(t)
	x := (string)(e.lhs)
	t[x] = v
	if v != TyIllTyped {
		return true, DECL, vP
	}
	return false, DECL, vP
}

// Check assign

func (assign Assign) check(t TyState) (bool, int, int) {
	v, vP := assign.value.infer(t)
	x, ok := t[assign.name]
	if ok {
		if v != TyIllTyped && x == v {
			x = v
			return true, ASSIGN, vP
		} else {
			return false, ASSIGN, vP
		}
	} else {
		return false, ASSIGN, Variables
	}
	return false, ASSIGN, vP
}

// Check coms

func (e ComS) check(t TyState) (bool, int, int) {
	v, vP, vPi := e[0].check(t)
	x, xP, xPi := e[1].check(t)
	if v && x {
		return true, COMS, 0
	}
	if !v {
		return false, vP, vPi
	} else {
		return false, xP, xPi
	}
}

func (e Print) check(t TyState) (bool, int, int) {
	v, vPi := e.e.infer(t)
	if v != TyIllTyped {
		return true, PRINT, vPi
	}
	return false, PRINT, vPi
}

// Block

func (b Block) check(t TyState) (bool, int, int) {
	v, vP, vPi := b.s.check(t)
	if v {
		return true, BLOCK, vPi
	} else {
		return false, vP, vPi
	}
}

// If Else

func (ifel IfEl) check(t TyState) (bool, int, int) {
	b1, _ := ifel.e.infer(t)
	b2, b2P, _ := ifel.b1.check(t)
	b3, b3P, _ := ifel.b2.check(t)
	if b1 == TyBool && b2 && b3 {
		return true, IF, 0
	}
	if b1 != TyBool {
		return false, IF, Condition
	}
	if !b2 {
		return false, b2P, BlockT
	} else {
		return false, b3P, BlockT
	}

}

// While

func (w While) check(t TyState) (bool, int, int) {
	b1, _ := w.e.infer(t)
	b2, _, _ := w.b.check(t)
	if b1 == TyBool && b2 {
		return true, WHILE, 0
	}
	if b1 != TyBool {
		return false, WHILE, Condition
	} else {
		return false, WHILE, BlockT
	}

}

// Simple scanner/lexer

// Tokens
const (
	EOS    = 0
	ZERO   = 1
	ONE    = 2
	TWO    = 3
	OPEN   = 4
	CLOSE  = 5
	PLUS   = 6
	MULT   = 7
	THREE  = 8
	FOUR   = 9
	FIVE   = 10
	SIX    = 11
	SEVEN  = 12
	EIGHT  = 13
	NINE   = 14
	LESS   = 15
	COMS   = 16
	EQU    = 17
	AND    = 18
	OR     = 19
	TRUE   = 20
	FALSE  = 21
	NEG    = 22
	VAR    = 23
	ASSIGN = 24
	DECL   = 25
	WHILE  = 26
	IF     = 27
	PRINT  = 28
	OPENC  = 29
	CLOSEC = 30
	ELSE   = 31
	BLOCK  = 32
)

func (s State) printToken() string {
	switch {
	case s.tok == 0:
		return "EOS"
	case s.tok == 1:
		return "ZERO"
	case s.tok == 2:
		return "ONE"
	case s.tok == 3:
		return "TWO"
	case s.tok == 4:
		return "OPEN"
	case s.tok == 5:
		return "CLOSE"
	case s.tok == 6:
		return "PLUS"
	case s.tok == 7:
		return "MULT"
	case s.tok == 8:
		return "THREE"
	case s.tok == 9:
		return "FOUR"
	case s.tok == 10:
		return "FIVE"
	case s.tok == 11:
		return "SIX"
	case s.tok == 12:
		return "SEVEN"
	case s.tok == 13:
		return "EIGHT"
	case s.tok == 14:
		return "NINE"
	case s.tok == 15:
		return "LESS"
	case s.tok == 16:
		return "COMS"
	case s.tok == 17:
		return "EQU"
	case s.tok == 18:
		return "AND"
	case s.tok == 19:
		return "OR"
	case s.tok == 20:
		return "True"
	case s.tok == 21:
		return "FALSE"
	case s.tok == 22:
		return "NEG"
	case s.tok == 23:
		return "VAR"
	case s.tok == 24:
		return "ASSIGN"
	case s.tok == 25:
		return "DECL"
	case s.tok == 26:
		return "WHILE"
	case s.tok == 27:
		return "IF"
	case s.tok == 28:
		return "PRINT"
	case s.tok == 29:
		return "OPENC"
	case s.tok == 30:
		return "CLOSEC"
	case s.tok == 31:
		return "ELSE"

	}
	return "Not a Token"
}

func printToken(i int) string {
	switch {
	case i == 0:
		return "EOS"
	case i == 1:
		return "ZERO"
	case i == 2:
		return "ONE"
	case i == 3:
		return "TWO"
	case i == 4:
		return "OPEN"
	case i == 5:
		return "CLOSE"
	case i == 6:
		return "PLUS"
	case i == 7:
		return "MULT"
	case i == 8:
		return "THREE"
	case i == 9:
		return "FOUR"
	case i == 10:
		return "FIVE"
	case i == 11:
		return "SIX"
	case i == 12:
		return "SEVEN"
	case i == 13:
		return "EIGHT"
	case i == 14:
		return "NINE"
	case i == 15:
		return "LESS"
	case i == 16:
		return "COMS"
	case i == 17:
		return "EQU"
	case i == 18:
		return "AND"
	case i == 19:
		return "OR"
	case i == 20:
		return "True"
	case i == 21:
		return "FALSE"
	case i == 22:
		return "NEG"
	case i == 23:
		return "VAR"
	case i == 24:
		return "ASSIGN"
	case i == 25:
		return "DECL"
	case i == 26:
		return "WHILE"
	case i == 27:
		return "IF"
	case i == 28:
		return "PRINT"
	case i == 29:
		return "OPENC"
	case i == 30:
		return "CLOSEC"
	case i == 31:
		return "ELSE"
	case i == 32:
		return "BLOCK"
	}
	return "Not a Token"
}

func printExp(i int) string {
	switch {
	case i == 1:
		return "Integer"
	case i == 2:
		return "Booleans"
	case i == 3:
		return "IllTyped Addition"
	case i == 4:
		return "IllTyped Multiplication"
	case i == 5:
		return "IllTyped Disjunction"
	case i == 6:
		return "IllTyped Conjuction"
	case i == 7:
		return "IllTyped Negation"
	case i == 8:
		return "IllTyped Equality"
	case i == 9:
		return "IllTyped Lesser"
	case i == 10:
		return "Variable not declarated"
	case i == 11:
		return "Condition IllTyped"
	case i == 12:
		return "Error in Block "
	default:
		return "Undefined"
	}
}

func scan(s string) (string, int) {
	for {
		switch {
		case len(s) == 0:
			return s, EOS
		case s[0] == '0':
			return s[1:len(s)], ZERO
		case s[0] == '1':
			return s[1:len(s)], ONE
		case s[0] == '2':
			return s[1:len(s)], TWO
		case s[0] == '3':
			return s[1:len(s)], THREE
		case s[0] == '4':
			return s[1:len(s)], FOUR
		case s[0] == '5':
			return s[1:len(s)], FIVE
		case s[0] == '6':
			return s[1:len(s)], SIX
		case s[0] == '7':
			return s[1:len(s)], SEVEN
		case s[0] == '8':
			return s[1:len(s)], EIGHT
		case s[0] == '9':
			return s[1:len(s)], NINE
		case s[0] == '+':
			return s[1:len(s)], PLUS
		case s[0] == '*':
			return s[1:len(s)], MULT
		case s[0] == '(':
			return s[1:len(s)], OPEN
		case s[0] == ')':
			return s[1:len(s)], CLOSE
		case s[0] == '{':
			return s[1:len(s)], OPENC
		case s[0] == '}':
			return s[1:len(s)], CLOSEC
		case s[0] == '<':
			return s[1:len(s)], LESS
		case s[0] == '!':
			return s[1:len(s)], NEG
		case len(s) >= 2 && s[0] == '=' && s[1] == '=':
			return s[2:len(s)], EQU
		case s[0] == '=':
			return s[1:len(s)], ASSIGN
		case len(s) >= 2 && s[0] == ':' && s[1] == '=':
			return s[2:len(s)], DECL
		case len(s) >= 2 && s[0] == '|' && s[1] == '|':
			return s[2:len(s)], OR
		case len(s) >= 2 && s[0] == '&' && s[1] == '&':
			return s[2:len(s)], AND
		case s[0] == ';':
			return s[1:len(s)], COMS
		case len(s) >= 2 && unicode.IsLetter(rune(s[0])):
			i := 0
			for len(s) >= i+1 && unicode.IsLetter(rune(s[0+i])) {
				i++
			}
			switch {
			case s[0:i] == "if":
				return s[i:len(s)], IF
			case s[0:i] == "else":
				return s[i:len(s)], ELSE
			case s[0:i] == "while":
				return s[i:len(s)], WHILE
			case s[0:i] == "print":
				return s[i:len(s)], PRINT
			case s[0:i] == "true":
				return s[i:len(s)], TRUE
			case s[0:i] == "false":
				return s[i:len(s)], FALSE
			default:
				varName = s[0:i]
				return s[i:len(s)], VAR
			}
		default:
			s = s[1:len(s)]
		}

	}
}

type State struct {
	s   *string
	tok int
}

func next(s *State) {
	s2, tok := scan(*s.s)

	s.s = &s2
	s.tok = tok
}

// Parse
func parseBlock(s *State) (bool, Block) {

	if s.tok != OPENC {

		return false, Block{}
	}

	b, t := parseComS(s)
	if !b {
		return false, Block{}
	}

	if !b || s.tok != CLOSEC {
		return false, Block{}
	}
	next(s)

	return true, Block{t}
}

func parseComS(s *State) (bool, Stmt) {
	b, e := parseStatement(s)
	if !b {
		return false, e
	}
	return parseComS2(s, e)
}

// Or2 ::= == T Or2
func parseComS2(s *State, e Stmt) (bool, Stmt) {
	if s.tok == COMS {

		b, f := parseStatement(s)
		if !b {
			return false, e
		}
		t := (ComS)([2]Stmt{e, f})

		return parseComS2(s, t)
	}

	return true, e
}

func parseStatement(s *State) (bool, Stmt) {
	next(s)

	switch {
	case s.tok == VAR:
		next(s)
		name := varName
		switch {
		case s.tok == DECL:
			next(s)
			b, e := parseOr(s)

			if !b {
				return false, Decl{}
			}
			return true, Decl{name, e}
		case s.tok == ASSIGN:
			next(s)
			b, e := parseOr(s)
			if !b {
				return false, Assign{}
			}
			return true, Assign{name, e}
		}
		return false, nil

	case s.tok == WHILE:

		next(s)
		b, e := parseOr(s)
		if !b {
			return false, While{}
		}

		b, bl := parseBlock(s)

		if !b {
			return false, While{}
		}
		return true, While{e, bl}

	case s.tok == IF:
		next(s)
		b, e := parseOr(s)

		if !b {
			return false, IfEl{}
		}

		b, bl := parseBlock(s)

		if !b {
			return false, IfEl{}
		}

		if s.tok != ELSE {
			return false, IfEl{}
		}
		next(s)
		b, bl2 := parseBlock(s)
		if !b {
			return false, IfEl{}
		}
		return true, IfEl{e, bl, bl2}
	case s.tok == PRINT:
		next(s)
		b, e := parseOr(s)
		if !b {
			return false, Print{}
		}
		return true, Print{e}
	default:
		return false, nil
	}

}

// Or ::= Equ And(2
func parseOr(s *State) (bool, Exp) {
	b, e := parseAnd(s)
	if !b {
		return false, e
	}
	return parseOr2(s, e)
}

// Or2 ::= == T Or2
func parseOr2(s *State, e Exp) (bool, Exp) {
	if s.tok == OR {
		next(s)
		b, f := parseAnd(s)
		if !b {
			return false, e
		}
		t := (Or)([2]Exp{e, f})
		return parseOr2(s, t)
	}

	return true, e
}

// And ::= Equ And(2
func parseAnd(s *State) (bool, Exp) {
	b, e := parseEqu(s)
	if !b {
		return false, e
	}
	return parseAnd2(s, e)
}

// And2 ::= == T And2
func parseAnd2(s *State, e Exp) (bool, Exp) {
	if s.tok == AND {
		next(s)
		b, f := parseEqu(s)
		if !b {
			return false, e
		}
		t := (And)([2]Exp{e, f})
		return parseAnd2(s, t)
	}

	return true, e
}

// EQU ::= L EQU2
func parseEqu(s *State) (bool, Exp) {
	b, e := parseNeg(s)
	if !b {
		return false, e
	}
	return parseEqu2(s, e)
}

// EQU2 ::= == T EQU2
func parseEqu2(s *State, e Exp) (bool, Exp) {
	if s.tok == EQU {
		next(s)
		b, f := parseNeg(s)
		if !b {
			return false, e
		}
		t := (Equ)([2]Exp{e, f})
		return parseEqu2(s, t)
	}

	return true, e
}

// Neg ::= Or Neg2
func parseNeg(s *State) (bool, Exp) {
	b, e := parseL(s)
	if !b {
		return false, e
	}
	return parseNeg2(s, e)
}

// Neg2 ::= == Or Neg2
func parseNeg2(s *State, e Exp) (bool, Exp) {
	if s.tok == NEG {
		next(s)
		b, f := parseL(s)
		if !b {
			return false, e
		}
		t := (Neg)([1]Exp{f})
		return parseNeg2(s, t)
	}

	return true, e
}

// L ::= E L2
func parseL(s *State) (bool, Exp) {
	b, e := parseE(s)
	if !b {
		return false, e
	}
	return parseL2(s, e)
}

// L2 ::= < T L2
func parseL2(s *State, e Exp) (bool, Exp) {
	if s.tok == LESS {
		next(s)
		b, f := parseE(s)
		if !b {
			return false, e
		}
		t := (Les)([2]Exp{e, f})
		return parseL2(s, t)
	}

	return true, e
}

// E  ::= T E2
func parseE(s *State) (bool, Exp) {
	b, e := parseT(s)
	if !b {
		return false, e
	}
	return parseE2(s, e)
}

// E2 ::= + T E2 |
func parseE2(s *State, e Exp) (bool, Exp) {
	if s.tok == PLUS {
		next(s)
		b, f := parseT(s)
		if !b {
			return false, e
		}
		t := (Plus)([2]Exp{e, f})
		return parseE2(s, t)
	}

	return true, e
}

// T  ::= F T2
func parseT(s *State) (bool, Exp) {
	b, e := parseF(s)
	if !b {
		return false, e
	}
	return parseT2(s, e)
}

// T2 ::= * F T2 |
func parseT2(s *State, e Exp) (bool, Exp) {
	if s.tok == MULT {
		next(s)
		b, f := parseF(s)
		if !b {
			return false, e
		}
		t := (Mult)([2]Exp{e, f})
		return parseT2(s, t)
	}
	return true, e
}

// F ::= N | (E)
func parseF(s *State) (bool, Exp) {
	switch {
	case s.tok == ZERO:
		next(s)
		return true, (Num)(0)
	case s.tok == ONE:
		next(s)
		return true, (Num)(1)
	case s.tok == TWO:
		next(s)
		return true, (Num)(2)
	case s.tok == THREE:
		next(s)
		return true, (Num)(3)
	case s.tok == FOUR:
		next(s)
		return true, (Num)(4)
	case s.tok == FIVE:
		next(s)
		return true, (Num)(5)
	case s.tok == SIX:
		next(s)
		return true, (Num)(6)
	case s.tok == SEVEN:
		next(s)
		return true, (Num)(7)
	case s.tok == EIGHT:
		next(s)
		return true, (Num)(8)
	case s.tok == NINE:
		next(s)
		return true, (Num)(9)
	case s.tok == TRUE:
		next(s)
		return true, (Bool)(true)
	case s.tok == FALSE:
		next(s)
		return true, (Bool)(false)
	case s.tok == OPEN:
		next(s)
		b, e := parseOr(s)
		if !b {
			return false, e
		}
		if s.tok != CLOSE {
			return false, e
		}
		next(s)
		return true, e
	case s.tok == NEG:
		return true, (Num)(0)
	case s.tok == VAR:
		next(s)
		return true, (Var)(varName)
	case s.tok == WHILE:
		return true, (Num)(0)
	case s.tok == OPENC:
		return true, (Num)(0)
	case s.tok == CLOSEC:
		return true, (Num)(0)
	case s.tok == IF:
		return true, (Num)(0)
	case s.tok == ELSE:
		return true, (Num)(0)
	case s.tok == DECL:
		return true, (Num)(0)
	case s.tok == ASSIGN:
		return true, (Num)(0)
	case s.tok == PRINT:
		return true, (Num)(0)
	}

	return false, (Num)(0)
}

func parse(s string) (bool, int, Block) {
	st := State{&s, EOS}
	inputLength = len(s)
	next(&st)
	_, e := parseBlock(&st)
	if st.tok == EOS {
		return true, 0, e
	}
	errorLength = len(*st.s)
	errorAt := inputLength - errorLength
	return false, errorAt, Block{} // dummy value
}

func debug(s string) {
	fmt.Printf("%s", s)
}

func test(s string) {
	stmt, errorAtStmt, e := parse(s)
	var vals = make(ValState)
	var types = make(TyState)
	fmt.Printf("\n Input: %s", s)
	if !stmt {
		fmt.Printf("\n ERROR ON PARSE \n AT CHARACTER %d \n", errorAtStmt)
		return
	}
	fmt.Printf("\n Output Parse: %s", e.pretty())

	exp, errorIn, errorAtExp := e.check(types)
	fmt.Printf("\n Check: %t ", exp)
	if !exp {
		fmt.Printf("\n ERROR ON EVALUATION \n")
		fmt.Printf(" Illtyped Statement found, StatementType = " + printToken(errorIn) + ", Reason = " + printExp(errorAtExp) + "\n")
		return
	}
	fmt.Printf("\n Evalutaion: ")
	e.eval(vals)
	fmt.Printf("\n")
}

func testParserGood() {

	fmt.Printf("\n Test 1.1 - Declaration \n")
	test("{varX:=3}")
	fmt.Printf("\n Test 1.2 - False Declaration Example - Error at char 8\n")
	test("{varX:==3}")

	fmt.Printf("\n Test 2.1 - Command Sequence - two expressions \n")
	test("{varX:=3;varY:=4}")
	fmt.Printf("\n Test 2.2 - Command Sequence - three expressions \n")
	test("{varX:=3;varY:=4;varZ:=7}")
	fmt.Printf("\n Test 2.3 - False Command Sequence - Error at char 18\n")
	test("{varX:=3;varY:=4;;varZ:=7}")

	fmt.Printf("\n Test 3.1 - Print - print expression \n")
	test("{print false; print true}")
	fmt.Printf("\n Test 3.2 - False Print - print expression \n")
	test("{print fasle; print true}")

	fmt.Printf("\n Test 4.1 - Assignment - set varX 4\n")
	test("{varX:=3;print varX; varX=4; print varX}")
	fmt.Printf("\n Test 4.2 - False Assignment - not declarated\n")
	test("{varX=4; print varX}")

	fmt.Printf("\n Test 5.1 - Plus - print varX + varY \n")
	test("{varX:=3;varY:=4;print varX+varY}")
	fmt.Printf("\n Test 5.2 - False Plus - bool + int\n")
	test("{varX:=true;varY:=4;print varX+varY}")

	fmt.Printf("\n Test 6.1 - Mult - print varX * varY \n")
	test("{varX:=3;varY:=4;print varX*varY}")
	fmt.Printf("\n Test 6.2 - False Mult - bool * int\n")
	test("{varX:=true;varY:=4;print varX*varY}")

	fmt.Printf("\n Test 7.1 - Less - print varX < varY \n")
	test("{varX:=3;varY:=4;print varX<varY}")
	fmt.Printf("\n Test 7.2 - False Less - bool < int\n")
	test("{varX:=true;varY:=4;print varX<varY}")

	fmt.Printf("\n Test 8.1 - And - print varX && varY \n")
	test("{varX:=true;varY:=true;print varX&&varY}")
	test("{varX:=false;varY:=true;print varX&&varY}")
	fmt.Printf("\n Test 8.2 - False And - bool && int\n")
	test("{varX:=true;varY:=4;print varX&&varY}")

	fmt.Printf("\n Test 9.1 - Or - print varX || varY \n")
	test("{varX:=true;varY:=false;print varX||varY}")
	test("{varX:=false;varY:=false;print varX||varY}")
	fmt.Printf("\n Test 9.2 - False Or - bool || int\n")
	test("{varX:=true;varY:=4;print varX||varY}")

	fmt.Printf("\n Test 10.1 - Equality - print varX == varY \n")
	test("{varX:=true;varY:=false;print varX==varY}")
	test("{varX:=false;varY:=false;print varX==varY}")
	test("{varX:=1;varY:=1;print varX==varY}")
	test("{varX:=1;varY:=2;print varX==varY}")
	fmt.Printf("\n Test 10.2 - False Equality - bool == int\n")
	test("{varX:=true;varY:=4;print varX==varY}")

	fmt.Printf("\n Test 11.1 - Negation - print !varX \n")
	test("{varX:=true;print !varX}")
	test("{varX:=false;print !varX}")
	fmt.Printf("\n Test 11.2 - False Negation - IllTyped \n")
	test("{varX:=1;print !varX}")

	fmt.Printf("\n Test 12.1 - If Else - if varX { print true } else { print false } \n")
	test("{varX:=true;if varX {print true} else {print false}}")
	test("{varX:=false;if varX {print true} else {print false}}")
	fmt.Printf("\n Test 12.2 - If Else - if varX == 1 { print true } else { print false } \n")
	test("{varX:=1;if varX == 1 {print true} else {print false}}")
	test("{varX:=2;if varX == 2 {print true} else {print false}}")
	fmt.Printf("\n Test 12.3 - If Else - IllTyped \n")
	test("{varX:=1;if varX {print true} else {print false}}")

	fmt.Printf("\n Test 13.1 - While - varX:=1; while varX<4 { print varX; varX+1}  \n")
	test("{varX:=1;while varX<4 {print varX; varX = varX+1}}")
	fmt.Printf("\n Test 13.2 - While - varX:=true; while varX { print varX; varX=false}  \n")
	test("{varX:=true;while varX {print varX; varX = false}}")
	fmt.Printf("\n Test 13.3 - While - IllTyped \n")
	test("{varX:=1;while varX {print varX}}")

}

// Helper functions to build ASTs by hand

func number(x int) Exp {
	return Num(x)
}

func variable(x string) Exp {
	return Var(x)
}

func boolean(x bool) Exp {
	return Bool(x)
}

func plus(x, y Exp) Exp {
	return (Plus)([2]Exp{x, y})
}

func mult(x, y Exp) Exp {
	return (Mult)([2]Exp{x, y})
}

func and(x, y Exp) Exp {
	return (And)([2]Exp{x, y})
}

func or(x, y Exp) Exp {
	return (Or)([2]Exp{x, y})
}

// Negation

func neg(x Exp) Exp {
	return (Neg)([1]Exp{x})
}

// Equality Test
func equ(x, y Exp) Exp {
	return (Equ)([2]Exp{x, y})
}

// Lesser Test

func les(x, y Exp) Exp {
	return (Les)([2]Exp{x, y})
}

// Vars

// Command Sequence
func cs(x, y Stmt) Stmt {
	return (ComS)([2]Stmt{x, y})
}

// Variable declaration
func decl(x string, y Exp) Decl {
	return Decl{x, y}
}

// Variable assignment
func assign(x string, y Exp) Assign {
	return Assign{x, y}
}

// While
func while(e Exp, b Block) While {
	return While{e, b}
}

// If-then-else
func ifel(e Exp, b1 Block, b2 Block) IfEl {
	return IfEl{e, b1, b2}
}

// Print
func print(e Exp) Print {
	return Print{e}
}

// Block
func block(s Stmt) Block {
	return Block{s}
}
func examplesAST() {
	ast1 := block(cs(cs(cs(decl("trudy", number(3)), print(variable("trudy"))), cs(assign("trudy", plus(variable("trudy"), number(3))), print(variable("trudy")))), while(les(variable("trudy"), number(13)), block(ifel(les(variable("trudy"), number(11)), block(cs(print(variable("trudy")), assign("trudy", plus(variable("trudy"), number(1))))), block(assign("trudy", plus(variable("trudy"), number(1)))))))))
	var vals = make(ValState)
	var types = make(TyState)
	fmt.Printf("%s", ast1.pretty())
	ast1.check(types)
	ast1.eval(vals)
}

func main() {
	examplesAST()

	fmt.Printf("\n")
	testParserGood()
}
