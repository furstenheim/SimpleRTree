package SimpleRTree

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
	"sync"
)

func TestNode_ComputeDistances32(t *testing.T) {
	ns := []struct {
		Node32
		mind, maxd float32
	}{
		{
			Node32{
				BBox: [4]float32{
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
			Node32{
				BBox: [4]float32{
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
			Node32{
				BBox: [4]float32{
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
			Node32{
				BBox: [4]float32{
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
		mind, maxd := vectorComputeDistances32(n.BBox, 5, 5)
		assert.Equal(t, n.mind, mind)
		assert.Equal(t, n.maxd, maxd)
	}

}

func TestSimpleRTree_FindNearestPoint32(t *testing.T) {
	const size = 20
	points := make([]float32, size*2)
	for i := 0; i < 2*size; i++ {
		points[i] = rand.Float32()
	}
	fp := FlatPoints32(points)
	r := New32().Load(fp)
	for i := 0; i < 1000; i++ {
		x, y := rand.Float32(), rand.Float32()
		x1, y1, _, found := r.FindNearestPoint(x, y)
		x2, y2, _ := fp.linearClosestPoint32(x, y)
		assert.True(t, found, "We should always find nearest")
		assert.Equal(t, x1, x2)
		assert.Equal(t, y1, y2)
	}

}

func TestSimpleRTree_FindNearestPointWithinOutOfBBox32(t *testing.T) {
	const size = 20
	points := make([]float32, size*2)
	for i := 0; i < 2*size; i++ {
		points[i] = rand.Float32()
	}
	fp := FlatPoints32(points)
	r := New32().Load(fp)
	x, y := float32(5.), float32(5.)
	_, _, _, found := r.FindNearestPointWithin(x, y, 1)
	assert.False(t, found, "Closest point is not within distance")

}

func TestSimpleRTree_FindNearestPointWithinEmptyWithinBBox32(t *testing.T) {
	points := []float32{0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0}
	fp := FlatPoints32(points)
	r := New32().Load(fp)
	x, y := float32(0.5), float32(0.5)
	_, _, _, found := r.FindNearestPointWithin(x, y, 0.25)
	assert.False(t, found, "Closest point is not within distance")
}

func TestSimpleRTree_FindNearestPointBig32(t *testing.T) {
	const size = 20000
	points := make([]float32, size*2)
	for i := 0; i < 2*size; i++ {
		points[i] = rand.Float32()
	}
	fp := FlatPoints32(points)
	r := New32().Load(fp)
	rtreePool := &sync.Pool{}
	r2 := New32WithOptions(Options{
		RTreePool: rtreePool,
	})
	r2.Load(fp)
	r2.Destroy()
	// Check pooling works correctly
	r3 := New32WithOptions(Options{
		RTreePool: rtreePool,
	}).Load(fp)
	for i := 0; i < 1000; i++ {
		x, y := rand.Float32(), rand.Float32()
		x1, y1, _, found := r.FindNearestPoint(x, y)
		x3, y3, _, found := r3.FindNearestPoint(x, y)
		assert.True(t, found, "We should always find nearest")
		x2, y2, _ := fp.linearClosestPoint32(x, y)
		assert.Equal(t, x1, x2, "X coordinate")
		assert.Equal(t, y1, y2, "Y coordinate")
		assert.Equal(t, x3, x2, "X coordinate pooled")
		assert.Equal(t, y3, y2, "Y coordinate pooled")
	}

}

func TestSimpleRTree_FindNearestPointBigUnsafeMode32(t *testing.T) {
	const size = 20000
	points := make([]float32, size*2)
	for i := 0; i < 2*size; i++ {
		points[i] = rand.Float32()
	}
	fp := FlatPoints32(points)
	r := New32WithOptions(Options{UnsafeConcurrencyMode:true}).Load(fp)
	for i := 0; i < 1000; i++ {
		x, y := rand.Float32(), rand.Float32()
		x1, y1, _, found := r.FindNearestPoint(x, y)
		assert.True(t, found, "We should always find nearest")
		x2, y2, _ := fp.linearClosestPoint32(x, y)
		assert.Equal(t, x1, x2, "X coordinate")
		assert.Equal(t, y1, y2, "Y coordinate")
	}

}

func TestComputeSize32(t *testing.T) {
	testCases := []struct {
		len      int
		expected int
	}{
		{
			10,
			3,
		},
		{
			1000,
			129,
		},
		{
			11250,
			1277,
		},
	}
	for _, tc := range testCases {
		final := computeSize(tc.len)
		assert.True(t, tc.expected < final)
	}
}

func Benchmark_ComputeDistances32(b *testing.B) {
	size := 1000000
	points := make([]float32, size+10)
	for i := 0; i < size+10; i++ {
		points[i] = rand.Float32()
	}
	b.ResetTimer()
	i := 0
	for n := 0; n < b.N; n++ {
		i++
		i %= size
		bbox := newVectorBBox32(points[i], points[i+1], points[i+2], points[i+3])
		x, y := points[i+4], points[i+5]
		_, _ = computeDistances32(bbox, x, y)
	}

}

func Benchmark_VectorComputeDistances32(b *testing.B) {
	size := 1000000
	points := make([]float32, size+10)
	for i := 0; i < size+10; i++ {
		points[i] = rand.Float32()
	}
	b.ResetTimer()
	i := 0
	for n := 0; n < b.N; n++ {
		i++
		i %= size
		bbox := newVectorBBox32(points[i], points[i+1], points[i+2], points[i+3])
		x, y := points[i+4], points[i+5]
		_, _ = vectorComputeDistances32(bbox, x, y)
	}

}

func BenchmarkSimpleRTree_Load32(b *testing.B) {
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
			points := make([]float32, size*2)
			for i := 0; i < 2*size; i++ {
				points[i] = rand.Float32()
			}
			fp := FlatPoints32(points)

			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = New32().Load(fp)
			}
		})
	}
}
func BenchmarkSimpleRTree_LoadPooled32(b *testing.B) {
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
		pool := &sync.Pool{}
		b.Run(bm.name, func(b *testing.B) {
			size := bm.size
			points := make([]float32, size*2)
			for i := 0; i < 2*size; i++ {
				points[i] = rand.Float32()
			}
			fp := FlatPoints32(points)

			r0 := New32WithOptions(Options{RTreePool: pool}).Load(fp)
			r0.Destroy()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				r := New32WithOptions(Options{UnsafeConcurrencyMode: true, RTreePool: pool}).Load(fp)
				r.Destroy()
			}
		})
	}
}

func BenchmarkSimpleRTree_FindNearestPoint32(b *testing.B) {
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
			points := make([]float32, size*2)
			for i := 0; i < 2*size; i++ {
				points[i] = rand.Float32()
			}
			fp := FlatPoints32(points)
			r := New32WithOptions(Options{UnsafeConcurrencyMode: true}).Load(fp)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				x, y := rand.Float32(), rand.Float32()
				_, _, _, _ = r.FindNearestPoint(x, y)
			}
		})
	}
}

func BenchmarkSimpleRTree_FindNearestPointMemory32(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"1000", 1000},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			size := bm.size
			points := make([]float32, size*2)
			for i := 0; i < 2*size; i++ {
				points[i] = rand.Float32()
			}
			fp := FlatPoints32(points)
			r := New32WithOptions(Options{UnsafeConcurrencyMode: true}).Load(fp)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				x, y := rand.Float32(), rand.Float32()
				_, _, _, _ = r.FindNearestPoint(x, y)
			}
		})
	}
}
func BenchmarkSimpleRTree_LoadMemory32(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"1000", 1000},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			size := bm.size
			points := make([]float32, size*2)
			for i := 0; i < 2*size; i++ {
				points[i] = rand.Float32()
			}
			fp := FlatPoints32(points)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = New32().Load(fp)
			}
		})
	}
}

func (fp FlatPoints32) linearClosestPoint32(x, y float32) (x1, y1, d float32) {

	d = float32(math.Inf(1))
	for i := 0; i < fp.Len(); i++ {
		x2, y2 := fp.GetPointAt(i)
		if d1 := (x-x2) * (x-x2) + (y-y2)* (y-y2); d1 < d {
			d = d1
			x1 = x2
			y1 = y2
		}
	}
	return
}
