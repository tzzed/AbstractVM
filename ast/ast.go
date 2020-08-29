package ast

import (
	"avm/token"
	"bytes"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	return p.TokenLiteral()
}

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token    token.Token
	IntValue int32
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral returns string token literal
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

type ShortLiteral struct {
	Token      token.Token
	ShortValue int16
}

func (sl *ShortLiteral) expressionNode() {}

// TokenLiteral returns string token literal
func (sl *ShortLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *ShortLiteral) String() string {
	return sl.Token.Literal
}

type ByteLiteral struct {
	Token     token.Token
	ByteValue int8
}

func (bl *ByteLiteral) expressionNode() {}

// TokenLiteral returns string token literal
func (bl *ByteLiteral) TokenLiteral() string {
	return bl.Token.Literal
}

func (bl *ByteLiteral) String() string {
	return bl.Token.Literal
}

type FloatLiteral struct {
	Token      token.Token
	FloatValue float32
}

func (fl *FloatLiteral) expressionNode() {}

// TokenLiteral returns string token literal
func (fl *FloatLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FloatLiteral) String() string {
	return fl.Token.Literal
}

type DoubleLiteral struct {
	Token       token.Token
	DoubleValue float64
}

func (dl *DoubleLiteral) expressionNode() {}

// TokenLiteral returns string token literal
func (dl *DoubleLiteral) TokenLiteral() string {
	return dl.Token.Literal
}

func (dl *DoubleLiteral) String() string {
	return dl.Token.Literal
}

type InstructionStatement struct {
	Token token.Token
	Name  *Identifier
}

func (is *InstructionStatement) statementNode() {}

// TokenLiteral returns string token literal
func (is *InstructionStatement) TokenLiteral() string {
	return is.Token.Literal
}

func (is *InstructionStatement) String() string {
	var out bytes.Buffer

	out.WriteString(is.TokenLiteral() + " ")
	out.WriteString(is.Name.String())
	return out.String()
}

type PushStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *PushStatement) statementNode() {}

// TokenLiteral returns string token literal
func (ls *PushStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *PushStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(token.LPAREN)

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(token.RPAREN)

	return out.String()
}

type ValueLiteral struct {
	Token     token.Token
	Parameter *Identifier
}

func (ps *ValueLiteral) expressionNode() {}

// TokenLiteral returns string token literal.
func (ps *ValueLiteral) TokenLiteral() string {
	return ps.Token.Literal
}

func (ps *ValueLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(ps.TokenLiteral())
	out.WriteString(token.LPAREN + ps.Parameter.TokenLiteral())
	out.WriteString(token.RPAREN)

	return out.String()
}

type AssertStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (as *AssertStatement) statementNode() {}

func (as *AssertStatement) TokenLiteral() string {
	return as.Token.Literal
}

func (as *AssertStatement) String() string {
	var out bytes.Buffer

	out.WriteString(as.TokenLiteral() + " ")
	out.WriteString(as.Name.String())
	out.WriteString(token.LPAREN)

	if as.Value != nil {
		out.WriteString(as.Value.String())
	}

	out.WriteString(token.RPAREN)
	return out.String()
}

type AddStatement struct {
	Token token.Token
	Name  *Identifier
}

func (as *AddStatement) statementNode() {}

// TokenLiteral returns string token literal
func (as *AddStatement) TokenLiteral() string {
	return as.Token.Literal
}

func (as *AddStatement) String() string {
	var out bytes.Buffer

	out.WriteString(as.TokenLiteral())

	return out.String()
}

type PopStatement struct {
	Token token.Token
	Name  *Identifier
}

func (ps *PopStatement) statementNode() {}

// TokenLiteral returns string token literal.
func (ps *PopStatement) TokenLiteral() string {
	return ps.Token.Literal
}

func (ps *PopStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ps.TokenLiteral())

	return out.String()
}

type DivStatement struct {
	Token token.Token
	Name  *Identifier
}

func (do *DivStatement) statementNode() {}

func (do *DivStatement) TokenLiteral() string {
	return do.Token.Literal
}

func (do *DivStatement) String() string {
	var out bytes.Buffer

	out.WriteString(do.TokenLiteral())

	return out.String()
}

type MulStatement struct {
	Token token.Token
	Name  *Identifier
}

func (mo *MulStatement) statementNode() {}

// TokenLiteral returns string token literal
func (mo *MulStatement) TokenLiteral() string {
	return mo.Token.Literal
}

func (mo *MulStatement) String() string {
	var out bytes.Buffer

	out.WriteString(mo.TokenLiteral())

	return out.String()
}

type ModStatement struct {
	Token token.Token
	Name  *Identifier
}

func (mods *ModStatement) statementNode() {}

// TokenLiteral returns string token literal
func (mods *ModStatement) TokenLiteral() string {
	return mods.Token.Literal
}

func (mods *ModStatement) String() string {
	var out bytes.Buffer

	out.WriteString(mods.TokenLiteral())

	return out.String()
}
