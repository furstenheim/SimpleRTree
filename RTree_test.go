package SimpleRTree

import (
	"testing"
	"math/rand"
	"math"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert"
)

func TestSimpleRTree_FindNearestPoint(t *testing.T) {
	const size = 20
	points := make([]float64, size * 2)
	for i := 0; i < 2 * size; i++ {
		points[i] = rand.Float64()
	}
	fp := flatPoints(points)
	r := New().Load(fp)
	 for i := 0; i < 1000; i++ {
		x, y := rand.Float64(), rand.Float64()
		x1, y1 := r.FindNearestPoint(x, y)
		x2, y2 := fp.linearClosestPoint(x, y)
		assert.Equal(t, x1, x2)
		assert.Equal(t, y1, y2)
	}

}


type flatPoints []float64

func (fp flatPoints) Len () int {
	return len(fp) / 2
}

func (fp flatPoints) Swap (i, j int) {
	fp[2 * i], fp[2 * i + 1], fp[2 * j], fp[2 * j + 1] = fp[2 * j], fp[2 * j + 1], fp[2 * i], fp[2 * i + 1]
}

func (fp flatPoints) GetPointAt(i int) (x1, y1 float64) {
	return fp[2 * i], fp[2 * i +1]
}

func (fp flatPoints) linearClosestPoint (x, y float64) (x1, y1 float64) {
	d := math.Inf(1)
	for i := 0; i < fp.Len(); i++ {
		x2, y2 := fp.GetPointAt(i)
		if d1 := math.Pow(x - x2, 2) + math.Pow(y - y2, 2); d1 < d {
			d = d1
			x1 = x2
			y1 = y2
		}
	}
	return
}
