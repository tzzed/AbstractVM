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
	Token token.Token // the token.IDENT token
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
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *PushStatement) statementNode() {}

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
	Token token.Token // the token.LET token
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

type LoadStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LoadStatement) statementNode() {}

func (ls *LoadStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LoadStatement) String() string {
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

type StoreStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ss *StoreStatement) statementNode() {}

func (ss *StoreStatement) TokenLiteral() string {
	return ss.Token.Literal
}

func (ss *StoreStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ss.TokenLiteral() + " ")
	out.WriteString(ss.Name.String())
	out.WriteString(token.LPAREN)

	if ss.Value != nil {
		out.WriteString(ss.Value.String())
	}

	out.WriteString(token.RPAREN)
	return out.String()
}
