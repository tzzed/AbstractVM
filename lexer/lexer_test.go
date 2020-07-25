package lexer

import (
	"github.com/stretchr/testify/require"
	"testing"
)
import "avm/token"

func TestNextToken(t *testing.T) {
	input := `push int32(5)
			push int32(10)
			add
			push float(44.55)
			mul
			push double(42.42)
			push int32(42)

			dump
			pop
			assert double(42.42)
			exit
			;;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{
			token.PUSH, "push",
		},
		{
			token.INT32, "int32",
		},
		{
			token.LPAREN, "(",
		},
		{
			token.INT, "5",
		},
		{
			token.RPAREN, ")",
		},
		{
			token.PUSH, "push",
		},
		{
			token.INT32, "int32",
		},
		{
			token.LPAREN, "(",
		},
		{
			token.INT, "10",
		},
		{
			token.RPAREN, ")",
		},
		{
			token.ADD, "add",
		},
		{
			token.PUSH, "push",
		},
		{
			token.FLOAT, "float",
		},
		{
			token.LPAREN, "(",
		},
		{
			token.FLOAT_NUM, "44.55",
		},
		{
			token.RPAREN, ")",
		},
		{
			token.MUL, "mul",
		},

		{
			token.PUSH, "push",
		},
		{
			token.DOUBLE, "double",
		},
		{
			token.LPAREN, "(",
		},
		{
			token.FLOAT_NUM, "42.42",
		},
		{
			token.RPAREN, ")",
		},
		{
			token.PUSH, "push",
		},
		{
			token.INT32, "int32",
		},
		{
			token.LPAREN, "(",
		},
		{
			token.INT, "42",
		},
		{
			token.RPAREN, ")",
		},
		{
			token.DUMP, "dump",
		},
		{
			token.POP, "pop",
		},
		{
			token.ASSERT, "assert",
		},
		{
			token.DOUBLE, "double",
		},
		{
			token.LPAREN, "(",
		},
		{
			token.FLOAT_NUM, "42.42",
		},
		{
			token.RPAREN, ")",
		},
		{
			token.EXIT, "exit",
		},
		{
			token.EOI, ";;",
		},
		{
			token.EOF,
			"",
		},
	}

	l := New(input)
	for _, tt := range tests {
		tok := l.NextToken()
		require.Equal(t, tt.expectedType, tok.Type)
		require.Equal(t, tt.expectedLiteral, tok.Literal)
	}

}
