package evaluator

import (
	"avm/parser"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModOperand(t *testing.T) {
	tests := []struct {
		name  string
		input string
		a     Value
		b     Value
		want  Value
		fails bool
	}{
		{"mod / integer divide by 0 ", "mod", NewInt32Value(0), NewInt8Value(5), NewInt32Value(0), true},
		{"mod / short with result 0", "mod", NewInt16Value(5), NewInt8Value(2), NewInt16Value(2 % 5), false},
		{"mod / short ", "mod", NewInt16Value(8), NewInt8Value(32), NewInt16Value(32 % 8), false},
		{"mod / short divide by 0", "mod", NewInt16Value(8), NewInt8Value(0), NewInt16Value(0), false},
		{"mod / float ", "mod", NewFloatValue(3), NewFloatValue(32.33), NewFloatValue(float32(math.Mod(32.33, 3))), false},
	}

	for _, tt := range tests {
		t.Run("mod evaluator_"+tt.name, func(t *testing.T) {
			st := NewStack()
			st.Push(tt.a)
			st.Push(tt.b)
			ev, err := testEval(t, tt.input, st)
			if tt.fails {
				require.Error(t, err)
				return
			} else {
				require.Equal(t, tt.want.V, ev.V)
			}
		})

	}
}

func TestDivOperand(t *testing.T) {
	tests := []struct {
		input string
		a     Value
		b     Value
		want  Value
		fails bool
	}{
		{"div", NewInt32Value(0), NewInt8Value(5), NewInt32Value(0), true},
		{"div", NewInt16Value(5), NewInt8Value(2), NewInt16Value(0), false},
		{"div", NewInt16Value(8), NewInt8Value(32), NewInt16Value(4), false},
		{"div", NewInt16Value(8), NewInt8Value(0), NewInt16Value(0), false},
		{"div", NewFloatValue(3), NewFloatValue(32.33), NewFloatValue((32.33 / 3) + 0.000001), false},
	}

	for _, tt := range tests {
		t.Run("div evaluator", func(t *testing.T) {
			st := NewStack()
			st.Push(tt.a)
			st.Push(tt.b)
			ev, err := testEval(t, tt.input, st)
			if tt.fails {
				require.Error(t, err)
				return
			} else {
				require.Equal(t, tt.want.V, ev.V)
			}
		})

	}
}

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
		ev, err := testEval(t, tt.input, st)
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
		v, err := testEval(t, tt.input, st)
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
		v, err := testEval(t, tt.input, st)
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
		v, err := testEval(t, tt.input, st)
		require.NoError(t, err)
		require.Equal(t, tt.want.V, v.V)
	}

	require.Equal(t, len(tests), st.size)
}

func testEval(t *testing.T, input string, st *Stack) (Value, error) {
	p := parser.NewParser(input)
	pg, err := p.ParseInstruction()
	require.NoError(t, err)
	return st.Eval(pg)
}
