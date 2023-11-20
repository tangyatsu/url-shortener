package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 2",
			size: 2,
		},
		{
			name: "size = 8",
			size: 8,
		},
		{
			name: "size = 20",
			size: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str1 := NewRandomString(tt.size)
			str2 := NewRandomString(tt.size)

			assert.Len(t, str1, tt.size)
			assert.Len(t, str2, tt.size)

			assert.NotEqual(t, str1, str2) // at least we can check that we have generated two different strings
		})
	}
}
