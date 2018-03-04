package SimpleRTree

import (
	"testing"
	"math/rand"
	"math"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert"
	"fmt"
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
		x1, y1, found := r.FindNearestPoint(x, y)
		x2, y2 := fp.linearClosestPoint(x, y)
		assert.True(t, found, "We should always find nearest")
		assert.Equal(t, x1, x2)
		assert.Equal(t, y1, y2)
	}

}

func TestSimpleRTree_FindNearestPointBig(t *testing.T) {
	const size = 20000
	points := make([]float64, size * 2)
	for i := 0; i < 2 * size; i++ {
		points[i] = rand.Float64()
	}
	fp := flatPoints(points)
	r := New().Load(fp)
	fmt.Println("Finished loading")
	for i := 0; i < 1000; i++ {
		x, y := rand.Float64(), rand.Float64()
		x1, y1, found := r.FindNearestPoint(x, y)
		assert.True(t, found, "We should always find nearest")
		x2, y2 := fp.linearClosestPoint(x, y)
		assert.Equal(t, x1, x2)
		assert.Equal(t, y1, y2)
	}

}




func BenchmarkSimpleRTree_FindNearestPoint(b *testing.B) {
	benchmarks := []struct{
		name string
		size int
	}{
		{"10", 10,},
		{"1000", 1000,},
		{"10000", 10000,},
		{"100000", 100000,},
		{"200000", 200000,},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func (b * testing.B) {
			size := bm.size
			points := make([]float64, size * 2)
			for i := 0; i < 2 * size; i++ {
				points[i] = rand.Float64()
			}
			fp := flatPoints(points)
			r := New().Load(fp)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				x, y := rand.Float64(), rand.Float64()
				_, _, _ = r.FindNearestPoint(x, y)
			}
		})
	}
}

func BenchmarkSimpleRTree_FindNearestPointMemory(b *testing.B) {
	benchmarks := []struct{
		name string
		size int
	}{
		{"100", 100,},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func (b * testing.B) {
			size := bm.size
			points := make([]float64, size * 2)
			for i := 0; i < 2 * size; i++ {
				points[i] = rand.Float64()
			}
			fp := flatPoints(points)
			r := New().Load(fp)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				x, y := rand.Float64(), rand.Float64()
				_, _, _ = r.FindNearestPoint(x, y)
			}
		})
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
