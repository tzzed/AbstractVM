package evaluator

import (
	"avm/ast/object"
	"avm/lexer"
	"avm/parser"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvalPushExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"push int32(5)", int32(5)},
		{"push int16(5)", int16(5)},
		{"push int8(5)", int8(5)},
		{"push float(5.5)", float32(5.5)},
		{"push double(5.5)", 5.5},
		{"push float(5)", float32(5)},
	}

	st := New()
	for _, tt := range tests {
		evaluated := testEval(tt.input, st)

		require.True(t, testObjectType(t, evaluated, tt.expected))


	}

	require.Equal(t, len(tests), st.Stack.Len())
}

func testEval(input string, st *Stack) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	pg := p.ParseProgram()
	obj := st.Eval(pg)
	return obj
}




