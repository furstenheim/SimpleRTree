package SimpleRTree

type VectorBBox [4]float64

const (
	VECTOR_BBOX_MIN_X = 0
	VECTOR_BBOX_MIN_Y = 1
	VECTOR_BBOX_MAX_X = 2
	VECTOR_BBOX_MAX_Y = 3
)

func newVectorBBox(MinX, MinY, MaxX, MaxY float64) VectorBBox {
	return [4]float64{MinX, MinY, MaxX, MaxY}
}

func bbox2VectorBBox(b BBox) VectorBBox {
	return newVectorBBox(b.MinX, b.MinY, b.MaxX, b.MaxY)
}

/**
Code from
https://github.com/slimsag/rand/blob/master/simd/vec64.go
*/
func vectorBBoxExtend(b1, b2 VectorBBox) VectorBBox{
	return [4]float64{
		minFloat(b1[0], b2[0]),
		minFloat(b1[1], b2[1]),
		maxFloat(b1[2], b2[2]),
		maxFloat(b1[3], b2[3]),
	}
}

func (b1 VectorBBox) toBBox() BBox {
	return BBox{
		MinX: b1[VECTOR_BBOX_MIN_X],
		MinY: b1[VECTOR_BBOX_MIN_Y],
		MaxX: b1[VECTOR_BBOX_MAX_X],
		MaxY: b1[VECTOR_BBOX_MAX_Y],
	}
}
