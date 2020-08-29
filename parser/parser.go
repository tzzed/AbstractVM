package parser

import (
	"avm/ast"
	"avm/lexer"
	"avm/token"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	prefixParseFn func() (ast.Expression, error)
)

// Error returns the string representation of the error.
func (e *ParseError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s at char %d", e.Message, e.Pos)
	}
	return fmt.Sprintf("found %s, expected %s at char %d", e.Found, strings.Join(e.Expected, ", "), e.Pos)
}

func LookupOperand(op string) token.TokenType {
	if op, ok := operands[op]; ok {
		return op
	}

	return token.IDENT
}

type ParseError struct {
	Message  string
	Found    string
	Expected []string
	Pos      int
}

func newParseError(found string, expected []string, pos int) *ParseError {
	return &ParseError{Found: found,
		Expected: expected,
		Pos:      pos,
	}
}

type Parser struct {
	l              *lexer.Lexer
	curTok         token.Token
	peekTok        token.Token
	Errors         *ParseError
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

func (p *Parser) parseExit() {
	os.Exit(0)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	var pg ast.Program

	pg.Statements = []ast.Statement{}
	for p.curTok.Type != token.EOF && p.curTok.Type != token.EOI {
		if p.curTok.Type == token.EXIT {
			p.parseExit()
		}
		if p.curTok.Type == token.SEMICOLON {
			_, _ = p.parseCommentStatement()
			return nil, nil
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		pg.Statements = append(pg.Statements, stmt)
		p.nextToken()
	}

	return &pg, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.curTok.Type {
	case token.PUSH:
		return p.parsePushStatement()
	case token.ASSERT:
		return p.parseAssertStatement()
	case token.ADD:
		return p.parseAddStatement()
	case token.MUL:
		return p.parseMulStatement()
	case token.DIV:
		return p.parseDivStatement()
	case token.SEMICOLON:
		return p.parseCommentStatement()
	case token.POP:
		return p.parsePopStatement()
	case token.MOD:
		return p.parseModStatement()

	default:
		if token.IsIdent(p.curTok.Literal) {
			return p.parseInstructionStatement()
		}

	}

	return nil, nil
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

func (p *Parser) parseExpression() (ast.Expression, error) {
	prefix := p.prefixParseFns[p.curTok.Type]
	if prefix == nil {
		return nil, newParseError(p.curTok.Literal, []string{"identifier"}, p.l.Pos)
	}

	return prefix()
}

// parseInstructionStatement
func (p *Parser) parseInstructionStatement() (*ast.InstructionStatement, error) {
	stmt := &ast.InstructionStatement{Token: p.curTok}

	if !token.IsIdent(p.curTok.Literal) {
		return nil, newParseError(p.curTok.Literal, token.GetAllInstructions(), p.l.Pos)
	}

	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	return stmt, nil
}

func (p *Parser) parsePushStatement() (*ast.PushStatement, error) {
	stmt := &ast.PushStatement{Token: p.curTok}
	p.nextToken()
	operand := LookupOperand(p.curTok.Literal)
	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	if !p.expectPeek(token.LPAREN) {
		return nil, newParseError(p.curTok.Literal, token.GetAllOperands(), p.l.Pos)
	}

	p.nextToken()
	p.curTok.Type = operand
	stmt.Value, _ = p.parseExpression()
	return stmt, nil
}

func (p *Parser) parseIntegerLiteral() (ast.Expression, error) {
	lit := &ast.IntegerLiteral{Token: p.curTok}
	value, err := strconv.ParseInt(p.curTok.Literal, 0, 32)
	if err != nil {
		return nil, newParseError(p.curTok.Literal, []string{"identifier"}, p.l.Pos)
	}

	lit.IntValue = int32(value)
	return lit, nil
}

func (p *Parser) parseShortValueLiteral() (ast.Expression, error) {
	lit := &ast.ShortLiteral{Token: p.curTok}
	value, err := strconv.ParseInt(p.curTok.Literal, 0, 16)
	if err != nil {
		return nil, newParseError(p.curTok.Literal, []string{"identifier"}, p.l.Pos)
	}

	lit.ShortValue = int16(value)
	return lit, nil
}

func (p *Parser) parseByteValueLiteral() (ast.Expression, error) {
	lit := &ast.ByteLiteral{Token: p.curTok}
	value, err := strconv.ParseInt(p.curTok.Literal, 0, 8)
	if err != nil {
		return nil, newParseError(fmt.Sprintf("expected %s", token.INT8), []string{"value"}, p.l.Pos)
	}

	lit.ByteValue = int8(value)
	return lit, nil
}

func (p *Parser) parseFloatLiteral() (ast.Expression, error) {
	lit := &ast.FloatLiteral{Token: p.curTok}

	p.curTok.Type = token.FLOAT32
	value, err := strconv.ParseFloat(p.curTok.Literal, 32)
	if err != nil {
		return nil, newParseError(fmt.Sprintf("%s", token.INT8), []string{"value"}, p.l.Pos)
	}

	lit.FloatValue = float32(value)
	return lit, nil
}

func (p *Parser) parseDoubleLiteral() (ast.Expression, error) {
	lit := &ast.DoubleLiteral{Token: p.curTok}

	p.curTok.Type = token.FLOAT64
	value, err := strconv.ParseFloat(p.curTok.Literal, 64)
	if err != nil {
		return nil, newParseError(fmt.Sprintf("%s", p.curTok.Literal), []string{"value"}, p.l.Pos)
	}

	lit.DoubleValue = value
	return lit, nil
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

func (p *Parser) parseAssertStatement() (*ast.AssertStatement, error) {
	stmt := &ast.AssertStatement{Token: p.curTok}

	p.nextToken()
	operand := LookupOperand(p.curTok.Literal)
	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	if !p.expectPeek(token.LPAREN) {
		return nil, nil
	}

	p.nextToken()
	p.curTok.Type = operand
	stmt.Value, _ = p.parseExpression()
	return stmt, nil
}

func (p *Parser) parseAddStatement() (*ast.AddStatement, error) {
	stmt := &ast.AddStatement{Token: p.curTok}

	p.nextToken()
	operand := LookupOperand(p.curTok.Literal)
	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	if !p.expectPeek(token.LPAREN) {
		return nil, nil
	}

	p.nextToken()
	p.curTok.Type = operand
	return stmt, nil
}

func (p *Parser) parseCommentStatement() (*ast.InstructionStatement, error) {
	for {
		if token.EOF == p.curTok.Type || p.curTok.Type == token.LF {
			break
		}

		p.nextToken()
	}

	return nil, nil
}

func (p *Parser) parsePopStatement() (*ast.PopStatement, error) {
	stmt := &ast.PopStatement{Token: p.curTok}
	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	p.nextToken()
	if p.curTok.Literal != "" {
		return stmt, newParseError(p.curTok.Literal, []string{"end of instruction"}, p.l.Pos)
	}

	return stmt, nil
}

func (p *Parser) parseDivStatement() (*ast.DivStatement, error) {
	stmt := &ast.DivStatement{Token: p.curTok}
	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	p.nextToken()
	if p.curTok.Literal != "" {
		return stmt, newParseError(p.curTok.Literal, []string{"end of instruction"}, p.l.Pos)
	}

	return stmt, nil
}

func (p *Parser) parseMulStatement() (*ast.MulStatement, error) {
	stmt := &ast.MulStatement{Token: p.curTok}
	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	p.nextToken()
	if p.curTok.Literal != "" {
		return stmt, newParseError(p.curTok.Literal, []string{"end of instruction"}, p.l.Pos)
	}

	return stmt, nil
}

func (p *Parser) parseModStatement() (*ast.ModStatement, error) {
	stmt := &ast.ModStatement{Token: p.curTok}
	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	p.nextToken()
	if p.curTok.Literal != "" {
		return stmt, newParseError(p.curTok.Literal, []string{"end of instruction"}, p.l.Pos)
	}

	return stmt, nil
}
