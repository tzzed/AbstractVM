package parser

import (
	"avm/ast"
	"avm/lexer"
	"avm/token"
	"fmt"
	"strconv"
)

const (
	// Operands
	INT8       = "int8"
	INT16      = "int16"
	INT32      = "int32"
	FLOAT      = "float"
	DOUBLE     = "double"
	BIGDECIMAL = "bigdecimal"
	UNKNOWN    = "unknown"
)

var operands = map[string]token.TokenType{
	"int8":       INT8,
	"int16":      INT16,
	"int32":      INT32,
	"float":      FLOAT,
	"double":     DOUBLE,
	"bigdecimal": BIGDECIMAL,
	"unknown":    UNKNOWN,
}

type (
	prefixParseFn func() ast.Expression
)

func LookupOperand(op string) token.TokenType {
	if op, ok := operands[op]; ok {
		return op
	}

	return token.IDENT
}

type Parser struct {
	l              *lexer.Lexer
	curTok         token.Token
	peekTok        token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.INT8, p.parseByteValueLiteral)
	p.registerPrefix(token.INT16, p.parseShortValueLiteral)
	p.registerPrefix(token.INT32, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.DOUBLE, p.parseDoubleLiteral)

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	var pg ast.Program

	pg.Statements = []ast.Statement{}
	for p.curTok.Type != token.EOF && p.curTok.Type != token.EOI {
		stmt := p.parseStatement()
		pg.Statements = append(pg.Statements, stmt)
		p.nextToken()
	}

	return &pg
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curTok.Type {
	case token.PUSH:
		return p.parsePushStatement()
	case token.ASSERT:
		return p.parseAssertStatement()
	case token.LOAD:
		return p.parseLoadStatement()
	case token.STORE:
		return p.parseStoreStatement()
	default:
		if token.IsIdent(p.curTok.Literal) {
			return p.parseInstructionStatement()
		}

	}

	return nil
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curTok.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekTok.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	return false
}

func (p *Parser) parseExpression() ast.Expression {
	prefix := p.prefixParseFns[p.curTok.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curTok.Type)
		return nil
	}
	leftExp := prefix()

	return leftExp
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseInstructionStatement() *ast.InstructionStatement {
	stmt := &ast.InstructionStatement{Token: p.curTok}

	if !token.IsIdent(p.curTok.Literal) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	return stmt
}

func (p *Parser) parsePushStatement() *ast.PushStatement {
	stmt := &ast.PushStatement{Token: p.curTok}

	p.nextToken()
	operand := LookupOperand(p.curTok.Literal)
	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	p.curTok.Type = operand
	stmt.Value = p.parseExpression()
	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curTok}
	value, err := strconv.ParseInt(p.curTok.Literal, 0, 32)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curTok.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.IntValue = int32(value)

	return lit
}

func (p *Parser) parseShortValueLiteral() ast.Expression {
	lit := &ast.ShortLiteral{Token: p.curTok}
	value, err := strconv.ParseInt(p.curTok.Literal, 0, 16)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curTok.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.ShortValue = int16(value)

	return lit
}

func (p *Parser) parseByteValueLiteral() ast.Expression {
	lit := &ast.ByteLiteral{Token: p.curTok}
	value, err := strconv.ParseInt(p.curTok.Literal, 0, 8)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curTok.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.ByteValue = int8(value)

	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curTok}

	p.curTok.Type = token.FLOAT32
	value, err := strconv.ParseFloat(p.curTok.Literal, 32)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curTok.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.FloatValue = float32(value)
	return lit
}

func (p *Parser) parseDoubleLiteral() ast.Expression {
	lit := &ast.DoubleLiteral{Token: p.curTok}

	p.curTok.Type = token.FLOAT64
	value, err := strconv.ParseFloat(p.curTok.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curTok.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.DoubleValue = value
	return lit
}

func (p *Parser) parseValueLiteral() ast.Expression {
	lit := &ast.ValueLiteral{Token: p.curTok}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameter = p.parseValueParameter()

	return lit
}
func (p *Parser) parseValueParameter() *ast.Identifier {
	ident := &ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return ident
	}

	p.nextToken()

	ident = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return ident
}

func (p *Parser) parseAssertStatement() *ast.AssertStatement {
	stmt := &ast.AssertStatement{Token: p.curTok}

	p.nextToken()
	operand := LookupOperand(p.curTok.Literal)
	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	p.curTok.Type = operand
	stmt.Value = p.parseExpression()
	return stmt
}

func (p *Parser) parseLoadStatement() *ast.LoadStatement {
	stmt := &ast.LoadStatement{Token: p.curTok}

	p.nextToken()
	operand := LookupOperand(p.curTok.Literal)
	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	p.curTok.Type = operand
	stmt.Value = p.parseExpression()

	return stmt
}

func (p *Parser) parseStoreStatement() *ast.StoreStatement {
	stmt := &ast.StoreStatement{Token: p.curTok}

	p.nextToken()
	operand := LookupOperand(p.curTok.Literal)
	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	p.curTok.Type = operand
	stmt.Value = p.parseExpression()

	return stmt
}
