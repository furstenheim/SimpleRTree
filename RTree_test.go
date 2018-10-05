package SimpleRTree

import (
	"testing"
	"math/rand"
	"math"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert"
	"fmt"
	"log"
)

func TestNode_ComputeDistances (t *testing.T) {
	ns := []struct {
		Node
		string
	}{
		{
			Node{
				bbox: [4]float64{
					2,
					3,
					3,
					3,
				},
			},
			"greater equal",
		},
		{
			Node{
				bbox: [4]float64{
					1,
					3,
					2,
					3,
				},
			},
			"less equal",
		},		{
			Node{
				bbox: [4]float64{
					1,
					4,
					8,
					12,
				},
			},
			"less equal",
		},
	}
	for _, n := range(ns) {
		mind, maxd := vectorComputeDistances(n.bbox, 5, 5)
		log.Println(mind, maxd)
	}

}

func TestSimpleRTree_FindNearestPoint(t *testing.T) {
	const size = 20
	points := make([]float64, size * 2)
	for i := 0; i < 2 * size; i++ {
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
	points := make([]float64, size * 2)
	for i := 0; i < 2 * size; i++ {
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
	_, _, _, found := r.FindNearestPointWithin(x, y, 0.5)
	assert.False(t, found, "Closest point is not within distance")
}


func TestSimpleRTree_FindNearestPointBig(t *testing.T) {
	const size = 20000
	points := make([]float64, size * 2)
	for i := 0; i < 2 * size; i++ {
		points[i] = rand.Float64()
	}
	fp := FlatPoints(points)
	r := New().Load(fp)
	fmt.Println("Finished loading")
	for i := 0; i < 1000; i++ {
		x, y := rand.Float64(), rand.Float64()
		x1, y1, _, found := r.FindNearestPoint(x, y)
		assert.True(t, found, "We should always find nearest")
		x2, y2, _ := fp.linearClosestPoint(x, y)
		assert.Equal(t, x1, x2, "X coordinate")
		assert.Equal(t, y1, y2, "Y coordinate")
	}

}

func TestComputeSize (t *testing.T) {
	testCases := []struct{
		len int
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
	for _, tc := range(testCases) {
		final := computeSize(tc.len)
		assert.True(t, tc.expected < final)
	}
}


func BenchmarkSimpleRTree_Load(b *testing.B) {
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
			fp := FlatPoints(points)

			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = New().Load(fp)
			}
		})
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
			fp := FlatPoints(points)
			r := New().Load(fp)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				x, y := rand.Float64(), rand.Float64()
				_, _, _, _ = r.FindNearestPoint(x, y)
			}
		})
	}
}

func BenchmarkSimpleRTree_FindNearestPointMemory(b *testing.B) {
	benchmarks := []struct{
		name string
		size int
	}{
		{"1000", 1000,},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func (b * testing.B) {
			size := bm.size
			points := make([]float64, size * 2)
			for i := 0; i < 2 * size; i++ {
				points[i] = rand.Float64()
			}
			fp := FlatPoints(points)
			r := New().Load(fp)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				x, y := rand.Float64(), rand.Float64()
				_, _, _, _ = r.FindNearestPoint(x, y)
			}
		})
	}
}
func BenchmarkSimpleRTree_LoadMemory(b *testing.B) {
	benchmarks := []struct{
		name string
		size int
	}{
		{"1000", 1000,},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func (b * testing.B) {
			size := bm.size
			points := make([]float64, size * 2)
			for i := 0; i < 2 * size; i++ {
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


func (fp FlatPoints ) linearClosestPoint (x, y float64) (x1, y1, d1 float64) {
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
