package SimpleRTree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAvxBBox_Extend (t *testing.T) {
	testCases := []struct{
		b1, b2 BBox
		expected BBox
	} {
		{
			b1: BBox{0, 0, 1, 1},
			b2: BBox{1, 1, 2, 2},
			expected: BBox{0, 0, 2, 2},
		},
	}

	for _, tc := range(testCases) {
		aB1 := bbox2VectorBBox(tc.b1)
		aB2 := bbox2VectorBBox(tc.b2)
		vectorBBoxExtend(&aB1, &aB2)
		assert.Equal(t, tc.expected, aB1.toBBox())
		assert.Equal(t, tc.expected, tc.b1.extend(tc.b2))
	}
}
