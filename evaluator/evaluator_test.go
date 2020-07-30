package evaluator

import (
	"avm/lexer"
	"avm/parser"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMulOperand(t *testing.T) {
	tests := []struct {
		input   string
		a       Value
		b       Value
		want    Value
		errored bool
	}{
		{"mul", NewInt32Value(5), NewInt8Value(2), NewInt32Value(10), false},
		{"mul", NewInt16Value(5), NewInt8Value(2), NewInt16Value(10), false},
	}

	st := NewStack()

	for _, tt := range tests {
		st.Push(tt.a)
		st.Push(tt.b)
		ev, err := testEval(tt.input, st)
		if tt.errored {
			require.Error(t, err)
			continue
		} else {
			require.Equal(t, tt.want.V, ev.V)
		}
	}
}

func TestStackSwap(t *testing.T) {
	s := NewStack()
	s.Push(NewInt32Value(10))
	s.Push(NewInt32Value(14))
	err := s.Swap()
	require.NoError(t, err)
	a, err := s.Pop()
	require.NoError(t, err)
	require.Equal(t, NewInt32Value(10).V, a.V)
	b, err := s.Pop()
	require.NoError(t, err)
	require.Equal(t, NewInt32Value(14).V, b.V)
}

func TestStackDup(t *testing.T) {
	s := NewStack()
	s.Push(NewInt32Value(10))
	s.Dup()
	a, err := s.Pop()
	require.NoError(t, err)
	b, err := s.Pop()
	require.NoError(t, err)
	require.Equal(t, a.V, b.V)
}

func TestStackClear(t *testing.T) {
	s := NewStack()
	for i := 0; i < 10; i++ {
		s.Push(NewInt32Value(10))
	}

	s.Clear()
	require.Equal(t, 0, s.size)
}

func TestConvertAstToValue(t *testing.T) {
	tests := []struct {
		input   string
		a       Value
		want    Value
		errored bool
	}{
		{"assert int32(10)", NewInt32Value(10), NewInt32Value(10), false},
		{"assert int32(10)", NewFloatValue(10), NewInt32Value(10), true},
	}

	st := NewStack()
	for _, tt := range tests {
		st.Push(tt.a)
		v, err := testEval(tt.input, st)
		if tt.errored {
			require.Error(t, err)
			return
		}

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%.2f", tt.want.V), fmt.Sprintf("%.2f", v.V))
	}

	st.Dump()
}

func TestEvalAddInstruction(t *testing.T) {
	tests := []struct {
		input string
		a     Value
		b     Value
		want  Value
	}{
		{"add", NewInt8Value(5), NewInt8Value(5), NewInt8Value(10)},
		{"add", NewInt16Value(5), NewInt16Value(5), NewInt16Value(10)},
		{"add", NewInt16Value(5), NewInt32Value(5), NewInt32Value(10)},
		{"add", NewFloatValue(float32(42.42)), NewInt32Value(5), NewFloatValue(float32(47.42))},
		{"add", NewFloatValue(float32(42.42)), NewDoubleValue(42.42), NewDoubleValue(84.84)},
	}

	st := NewStack()
	for _, tt := range tests {
		st.Push(tt.a)
		st.Push(tt.b)
		v, err := testEval(tt.input, st)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%.2f", tt.want.V), fmt.Sprintf("%.2f", v.V))
	}

	require.Equal(t, len(tests), st.size)
	st.Dump()
}

func TestEvalPushInstruction(t *testing.T) {
	tests := []struct {
		input string
		want  Value
	}{
		{"push int32(5)", NewInt32Value(5)},
		{"push int16(5)", NewInt16Value(5)},
		{"push int8(5)", NewInt8Value(5)},
		{"push float(5.5)", NewFloatValue(5.5)},
		{"push double(5.5)", NewDoubleValue(5.5)},
		{"push float(5)", NewFloatValue(5)},
		{"push int32(5)", NewInt32Value(5)},
		{"push int16(5)", NewInt16Value(5)},
		{"push int8(5)", NewInt8Value(5)},
		{"push float(5.5)", NewFloatValue(5.5)},
		{"push double(5.5)", NewDoubleValue(5.5)},
	}

	st := NewStack()
	for _, tt := range tests {
		v, err := testEval(tt.input, st)
		require.NoError(t, err)
		require.Equal(t, tt.want.V, v.V)
	}

	require.Equal(t, len(tests), st.size)
}

func testEval(input string, st *Stack) (Value, error) {
	l := lexer.New(input)
	p := parser.New(l)
	pg := p.ParseProgram()
	return st.Eval(pg)
}
