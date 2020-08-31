package parser

import (
	"avm/ast"
	"avm/lexer"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestDumpStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		fails              bool
	}{
		{"dump", "dump", false},
		{"dump pop", "", true},
	}

	for _, tt := range tests {
		t.Run("dump statement", func(t *testing.T) {
			l := lexer.New(tt.input)
			p := NewParser(l)
			program, err := p.ParseInstruction()
			if tt.fails {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			stmt := program.Statements[0]
			testDumpStatement(t, stmt, tt.expectedIdentifier)

		})
	}
}

func TestModStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		fails              bool
	}{
		{"mod", "mod", false},
		{"mod pop", "", true},
	}

	for _, tt := range tests {
		t.Run("mod statement", func(t *testing.T) {
			l := lexer.New(tt.input)
			p := NewParser(l)
			program, err := p.ParseInstruction()
			if tt.fails {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			stmt := program.Statements[0]
			testModStatement(t, stmt, tt.expectedIdentifier)

		})
	}
}

func TestDivStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		fails              bool
	}{
		{"div", "div", false},
		{"div pop", "", true},
	}

	for _, tt := range tests {
		t.Run("div` statement", func(t *testing.T) {
			l := lexer.New(tt.input)
			p := NewParser(l)
			program, err := p.ParseInstruction()
			if tt.fails {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			stmt := program.Statements[0]
			testDivStatement(t, stmt, tt.expectedIdentifier)

		})
	}
}

func TestMulStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		fails              bool
	}{
		{"mul", "mul", false},
		{"mul pop", "", true},
	}

	for _, tt := range tests {
		t.Run("mul statement", func(t *testing.T) {
			l := lexer.New(tt.input)
			p := NewParser(l)
			program, err := p.ParseInstruction()
			if tt.fails {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			stmt := program.Statements[0]
			testMulStatement(t, stmt, tt.expectedIdentifier)

		})
	}

}

func TestPopStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		fails              bool
	}{
		{"pop", "pop", false},
		{"pop pop", "", true},
	}

	for _, tt := range tests {
		t.Run("pop statement", func(t *testing.T) {
			l := lexer.New(tt.input)
			p := NewParser(l)
			program, err := p.ParseInstruction()
			if tt.fails {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			stmt := program.Statements[0]
			testPopStatement(t, stmt, tt.expectedIdentifier)

		})
	}

}

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
		p := NewParser(l)
		program, _ := p.ParseInstruction()

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
		p := NewParser(l)
		stmt, _ := p.parseInstructionStatement()
		require.Equal(t, tt.want, stmt.TokenLiteral())
	}
}

func TestPushStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
		fail               bool
	}{
		{"push int8(42)", "int8", int8(42), false},
		{"push int16(42)", "int16", int16(42), false},
		{"push int32(42)", "int32", int32(42), false},
		{"push float(42.42)", "float", float32(42.42), false},
		{"push double(42.42)", "double", 42.42, false},
		{"push ", "double", 42.42, true},
	}

	for _, tt := range tests {
		t.Run("push statement", func(t *testing.T) {
			l := lexer.New(tt.input)
			p := NewParser(l)
			program, err := p.ParseInstruction()
			if tt.fail {
				require.Error(t, err)
				return
			}

			stmt := program.Statements[0]
			if !testPushStatement(t, stmt, tt.expectedIdentifier) {
				return
			}

			val := stmt.(*ast.PushStatement).Value
			if !testLiteralExpression(t, val, tt.expectedValue) {
				return
			}
		})

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

func testPopStatement(t *testing.T, s ast.Statement, name string) {
	require.Equal(t, "pop", s.TokenLiteral())
	popStmt, ok := s.(*ast.PopStatement)
	require.True(t, ok)
	require.Equal(t, name, popStmt.Name.TokenLiteral())
}

func testMulStatement(t *testing.T, s ast.Statement, name string) {
	require.Equal(t, "mul", s.TokenLiteral())
	mulStmt, ok := s.(*ast.MulStatement)
	require.True(t, ok)
	require.Equal(t, name, mulStmt.Name.TokenLiteral())
}

func testDivStatement(t *testing.T, s ast.Statement, name string) {
	require.Equal(t, "div", s.TokenLiteral())
	mulStmt, ok := s.(*ast.DivStatement)
	require.True(t, ok)
	require.Equal(t, name, mulStmt.Name.TokenLiteral())
}

func testModStatement(t *testing.T, s ast.Statement, name string) {
	require.Equal(t, "mod", s.TokenLiteral())
	mulStmt, ok := s.(*ast.ModStatement)
	require.True(t, ok)
	require.Equal(t, name, mulStmt.Name.TokenLiteral())
}

func testDumpStatement(t *testing.T, s ast.Statement, name string) {
	require.Equal(t, "dump", s.TokenLiteral())
	mulStmt, ok := s.(*ast.DumpStatement)
	require.True(t, ok)
	require.Equal(t, name, mulStmt.Name.TokenLiteral())
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
