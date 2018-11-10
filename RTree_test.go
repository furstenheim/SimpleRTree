package SimpleRTree

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
	"sync"
)

func TestNode_ComputeDistances(t *testing.T) {
	ns := []struct {
		Node
		mind, maxd float64
	}{
		{
			Node{
				BBox: [4]float64{
					2,
					3,
					3,
					3,
				},
			},
			8,
			8,
		},
		{
			Node{
				BBox: [4]float64{
					1,
					3,
					2,
					3,
				},
			},
			13,
			13,
		},
		{
			Node{
				BBox: [4]float64{
					1,
					4,
					8,
					12,
				},
			},
			0,
			17,
		},
		{
			Node{
				BBox: [4]float64{
					1,
					1,
					8,
					4,
				},
			},
			1,
			17,
		},
	}
	for _, n := range ns {
		mind, maxd := vectorComputeDistances(n.BBox, 5, 5)
		assert.Equal(t, n.mind, mind)
		assert.Equal(t, n.maxd, maxd)
	}

}

func TestSimpleRTree_FindNearestPoint(t *testing.T) {
	const size = 20
	points := make([]float64, size*2)
	for i := 0; i < 2*size; i++ {
		points[i] = rand.Float64()
	}
	fp := FlatPoints(points)
	r := New().Load(fp)
	for i := 0; i < 1000; i++ {
		x, y := rand.Float64(), rand.Float64()
		x1, y1, _, found := r.FindNearestPoint(x, y)
		x2, y2, _ := fp.linearClosestPoint(x, y)
		assert.True(t, found, "We should always find nearest")
		assert.Equal(t, x1, x2)
		assert.Equal(t, y1, y2)
	}

}

func TestSimpleRTree_FindNearestPointWithinOutOfBBox(t *testing.T) {
	const size = 20
	points := make([]float64, size*2)
	for i := 0; i < 2*size; i++ {
		points[i] = rand.Float64()
	}
	fp := FlatPoints(points)
	r := New().Load(fp)
	x, y := 5., 5.
	_, _, _, found := r.FindNearestPointWithin(x, y, 1)
	assert.False(t, found, "Closest point is not within distance")

}

func TestSimpleRTree_FindNearestPointWithinEmptyWithinBBox(t *testing.T) {
	points := []float64{0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0}
	fp := FlatPoints(points)
	r := New().Load(fp)
	x, y := 0.5, 0.5
	_, _, _, found := r.FindNearestPointWithin(x, y, 0.25)
	assert.False(t, found, "Closest point is not within distance")
}

func TestSimpleRTree_FindNearestPointBig(t *testing.T) {
	const size = 20000
	points := make([]float64, size*2)
	for i := 0; i < 2*size; i++ {
		points[i] = rand.Float64()
	}
	fp := FlatPoints(points)
	r := New().Load(fp)
	rtreePool := &sync.Pool{}
	r2 := NewWithOptions(Options{
		RTreePool: rtreePool,
	})
	r2.Load(fp)
	r2.Destroy()
	// Check pooling works correctly
	r3 := NewWithOptions(Options{
		RTreePool: rtreePool,
	}).Load(fp)
	for i := 0; i < 1000; i++ {
		x, y := rand.Float64(), rand.Float64()
		x1, y1, _, found := r.FindNearestPoint(x, y)
		x3, y3, _, found := r3.FindNearestPoint(x, y)
		assert.True(t, found, "We should always find nearest")
		x2, y2, _ := fp.linearClosestPoint(x, y)
		assert.Equal(t, x1, x2, "X coordinate")
		assert.Equal(t, y1, y2, "Y coordinate")
		assert.Equal(t, x3, x2, "X coordinate pooled")
		assert.Equal(t, y3, y2, "Y coordinate pooled")
	}

}

func TestSimpleRTree_FindNearestPointBigUnsafeMode(t *testing.T) {
	const size = 20000
	points := make([]float64, size*2)
	for i := 0; i < 2*size; i++ {
		points[i] = rand.Float64()
	}
	fp := FlatPoints(points)
	r := NewWithOptions(Options{UnsafeConcurrencyMode:true}).Load(fp)
	for i := 0; i < 1000; i++ {
		x, y := rand.Float64(), rand.Float64()
		x1, y1, _, found := r.FindNearestPoint(x, y)
		assert.True(t, found, "We should always find nearest")
		x2, y2, _ := fp.linearClosestPoint(x, y)
		assert.Equal(t, x1, x2, "X coordinate")
		assert.Equal(t, y1, y2, "Y coordinate")
	}

}

func TestComputeSize(t *testing.T) {
	testCases := []struct {
		len      int
		expected int
	}{
		{
			10,
			13,
		},
		{
			1000,
			1129,
		},
		{
			11250,
			11277,
		},
	}
	for _, tc := range testCases {
		final := computeSize(tc.len)
		assert.True(t, tc.expected < final)
	}
}

func Benchmark_ComputeDistances(b *testing.B) {
	size := 1000000
	points := make([]float64, size+10)
	for i := 0; i < size+10; i++ {
		points[i] = rand.Float64()
	}
	b.ResetTimer()
	i := 0
	for n := 0; n < b.N; n++ {
		i++
		i %= size
		bbox := newVectorBBox(points[i], points[i+1], points[i+2], points[i+3])
		x, y := points[i+4], points[i+5]
		_, _ = computeDistances(bbox, x, y)
	}

}

func Benchmark_VectorComputeDistances(b *testing.B) {
	size := 1000000
	points := make([]float64, size+10)
	for i := 0; i < size+10; i++ {
		points[i] = rand.Float64()
	}
	b.ResetTimer()
	i := 0
	for n := 0; n < b.N; n++ {
		i++
		i %= size
		bbox := newVectorBBox(points[i], points[i+1], points[i+2], points[i+3])
		x, y := points[i+4], points[i+5]
		_, _ = vectorComputeDistances(bbox, x, y)
	}

}

func BenchmarkSimpleRTree_Load(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"10", 10},
		{"1000", 1000},
		{"10000", 10000},
		{"100000", 100000},
		{"200000", 200000},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			size := bm.size
			points := make([]float64, size*2)
			for i := 0; i < 2*size; i++ {
				points[i] = rand.Float64()
			}
			fp := FlatPoints(points)

			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = New().Load(fp)
			}
		})
	}
}

func BenchmarkSimpleRTree_FindNearestPoint(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"10", 10},
		{"1000", 1000},
		{"10000", 10000},
		{"100000", 100000},
		{"200000", 200000},
		{"1000000", 1000000},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			size := bm.size
			points := make([]float64, size*2)
			for i := 0; i < 2*size; i++ {
				points[i] = rand.Float64()
			}
			fp := FlatPoints(points)
			r := NewWithOptions(Options{UnsafeConcurrencyMode: true}).Load(fp)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				x, y := rand.Float64(), rand.Float64()
				_, _, _, _ = r.FindNearestPoint(x, y)
			}
		})
	}
}

func BenchmarkSimpleRTree_FindNearestPointMemory(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"1000", 1000},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			size := bm.size
			points := make([]float64, size*2)
			for i := 0; i < 2*size; i++ {
				points[i] = rand.Float64()
			}
			fp := FlatPoints(points)
			r := NewWithOptions(Options{UnsafeConcurrencyMode: true}).Load(fp)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				x, y := rand.Float64(), rand.Float64()
				_, _, _, _ = r.FindNearestPoint(x, y)
			}
		})
	}
}
func BenchmarkSimpleRTree_LoadMemory(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"1000", 1000},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			size := bm.size
			points := make([]float64, size*2)
			for i := 0; i < 2*size; i++ {
				points[i] = rand.Float64()
			}
			fp := FlatPoints(points)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = New().Load(fp)
			}
		})
	}
}

func (fp FlatPoints) linearClosestPoint(x, y float64) (x1, y1, d1 float64) {
	d := math.Inf(1)
	for i := 0; i < fp.Len(); i++ {
		x2, y2 := fp.GetPointAt(i)
		if d1 := math.Pow(x-x2, 2) + math.Pow(y-y2, 2); d1 < d {
			d = d1
			x1 = x2
			y1 = y2
		}
	}
	return
}
