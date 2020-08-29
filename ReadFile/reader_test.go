package ReadFile

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadFile(t *testing.T) {
	tests := []struct{
		name string
		file string
		fails bool
	}{
		{"Good file format", "../f.avm", false},
		{"bad file format", "f.test", true},
		{"file not exist", "test.avm", true},
	}

	for _, tt := range tests {
		t.Run("files test", func(t *testing.T) {
			err := ReadFile(tt.file)
			if tt.fails {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
