package SimpleRTree

// Based on https://fgiesen.wordpress.com/2013/01/14/min-max-under-negation-and-an-aabb-trick/

type AvxBBox [4]float64

const (
	AVX_BBOX_MIN_X = 0
	AVX_BBOX_MIN_Y = 1
	AVX_BBOX_NEG_MAX_X = 2
	AVX_BBOX_NEG_MAX_Y = 3
)

func newAvxBBox (MinX, MinY, MaxX, MaxY float64) (AvxBBox){
	return [4]float64{MinX, MinY, -MaxX, -MaxY}
}

func bbox2AvxBBox (b BBox) (AvxBBox){
	return newAvxBBox(b.MinX, b.MinY, b.MaxX, b.MaxY)
}

/**
 Code from
 https://github.com/slimsag/rand/blob/master/simd/vec64.go
*/
// Implemented in avxBBox.s
func avxBBoxExtend(b1, b2 AvxBBox) AvxBBox

func (b1 AvxBBox) toBBox () BBox {
	return BBox{
		MinX: b1[AVX_BBOX_MIN_X],
		MinY: b1[AVX_BBOX_MIN_Y],
		MaxX: -b1[AVX_BBOX_NEG_MAX_X],
		MaxY: -b1[AVX_BBOX_NEG_MAX_Y],
	}
}
