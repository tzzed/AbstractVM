package evaluator

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetBiggerType(t *testing.T) {
	tests := []struct {
		a    Value
		b    Value
		want ValueType
	}{
		{NewInt8Value(1), NewInt16Value(1), ShortValue},
		{NewInt16Value(1), NewInt16Value(1), ShortValue},
		{NewInt16Value(1), NewInt32Value(1), IntegerValue},
		{NewInt16Value(1), NewFloatValue(1), FloatValue},
		{NewFloatValue(14.5), NewDoubleValue(42.42), DoubleValue},
	}

	for _, tt := range tests {
		require.Equal(t, GetBiggerType(tt.a, tt.b), tt.want)
	}
}
