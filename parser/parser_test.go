package parser

import (
	"avm/ast"
	"avm/lexer"
	"fmt"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestAssertStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"assert int8(42)", "int8", int8(42)},
		{"assert int16(42)", "int16", int16(42)},
		{"assert int32(42)", "int32", int32(42)},
		{"assert float(42.42)", "float", float32(42.42)},
		{"assert double(42.42)", "double", 42.42},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		stmt := program.Statements[0]
		testAssertStatement(t, stmt, tt.expectedIdentifier)

		val := stmt.(*ast.AssertStatement).Value
		testLiteralExpression(t, val, tt.expectedValue)
	}
}

func TestInstructionsStatement(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "pop",
			want:  "pop",
		},
		{
			"add",
			"add",
		},
		{
			"dump",
			"dump",
		},
		{
			"clear",
			"clear",
		},
		{
			"dup",
			"dup",
		},
		{
			"swap",
			"swap",
		},
		{
			"sub",
			"sub",
		},
		{
			"mul",
			"mul",
		},
		{
			"div",
			"div",
		},
		{
			"mod",
			"mod",
		},
		{
			"print",
			"print",
		},
		{
			"exit",
			"exit",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		stmt := p.parseInstructionStatement()
		require.Equal(t, tt.want, stmt.TokenLiteral())
	}
}

func TestPushStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"push int8(42)", "int8", int8(42)},
		{"push int16(42)", "int16", int16(42)},
		{"push int32(42)", "int32", int32(42)},
		{"push float(42.42)", "float", float32(42.42)},
		{"push double(42.42)", "double", 42.42},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		stmt := program.Statements[0]
		if !testPushStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.PushStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}

}
func testPushStatement(t *testing.T, s ast.Statement, name string) bool {
	require.Equal(t, "push", s.TokenLiteral())
	pushStmt, ok := s.(*ast.PushStatement)
	require.True(t, ok)
	require.Equal(t, name, pushStmt.Name.Value)
	require.Equal(t, name, pushStmt.Name.TokenLiteral())

	return true
}

func testAssertStatement(t *testing.T, s ast.Statement, name string) {
	require.Equal(t, "assert", s.TokenLiteral())
	asStmt, ok := s.(*ast.AssertStatement)
	require.True(t, ok)
	require.Equal(t, name, asStmt.Name.Value)
	require.Equal(t, name, asStmt.Name.TokenLiteral())
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {

	switch v := expected.(type) {
	case int8:
		return testByteLiteral(t, exp, v)
	case int16:
		return testShortLiteral(t, exp, v)
	case int32:
		return testIntegerLiteral(t, exp, v)
	case float32:
		return testFloatLiteral(t, exp, v)
	case float64:
		return testDoubleLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	default:
		fmt.Printf("%T\n", v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	require.True(t, ok)
	require.Equal(t, ident.Value, value)
	require.Equal(t, value, ident.TokenLiteral())

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, v int32) bool {
	intLiteral, ok := il.(*ast.IntegerLiteral)
	require.True(t, ok)

	require.Equal(t, intLiteral.IntValue, v)
	lit := strconv.Itoa(int(v))
	require.Equal(t, intLiteral.TokenLiteral(), lit)
	return true
}

func testShortLiteral(t *testing.T, il ast.Expression, v int16) bool {
	short, ok := il.(*ast.ShortLiteral)
	require.True(t, ok)

	require.Equal(t, short.ShortValue, v)
	lit := strconv.Itoa(int(v))
	require.Equal(t, short.TokenLiteral(), lit)
	return true
}

func testByteLiteral(t *testing.T, bl ast.Expression, v int8) bool {
	byteLiteral, ok := bl.(*ast.ByteLiteral)
	require.True(t, ok)

	require.Equal(t, byteLiteral.ByteValue, v)
	lit := strconv.Itoa(int(v))
	require.Equal(t, byteLiteral.TokenLiteral(), lit)
	return true
}

func testFloatLiteral(t *testing.T, fl ast.Expression, v float32) bool {
	floatLit, ok := fl.(*ast.FloatLiteral)
	require.True(t, ok)

	require.Equal(t, floatLit.FloatValue, v)
	lit := strconv.FormatFloat(float64(v), 'f', 2, 32)
	require.Equal(t, floatLit.TokenLiteral(), lit)
	return true
}

func testDoubleLiteral(t *testing.T, fl ast.Expression, v float64) bool {
	doubleLit, ok := fl.(*ast.DoubleLiteral)
	require.True(t, ok)

	require.Equal(t, doubleLit.DoubleValue, v)
	lit := strconv.FormatFloat(v, 'f', 2, 64)
	require.Equal(t, doubleLit.TokenLiteral(), lit)
	return true
}
