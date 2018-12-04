package SimpleRTree

type rVectorBBox [4]float64

const (
	vector_bbox_min_x = 0
	vector_bbox_min_y = 1
	vector_bbox_max_x = 2
	vector_bbox_max_y = 3
)

func newVectorBBox(MinX, MinY, MaxX, MaxY float64) rVectorBBox {
	return [4]float64{MinX, MinY, MaxX, MaxY}
}

func bbox2VectorBBox(b rBBox) rVectorBBox {
	return newVectorBBox(b.MinX, b.MinY, b.MaxX, b.MaxY)
}

/**
Code from
https://github.com/slimsag/rand/blob/master/simd/vec64.go
*/
func vectorBBoxExtend(b1, b2 rVectorBBox) rVectorBBox {
	return [4]float64{
		minFloat(b1[0], b2[0]),
		minFloat(b1[1], b2[1]),
		maxFloat(b1[2], b2[2]),
		maxFloat(b1[3], b2[3]),
	}
}

func (b1 rVectorBBox) toBBox() rBBox {
	return rBBox{
		MinX: b1[vector_bbox_min_x],
		MinY: b1[vector_bbox_min_y],
		MaxX: b1[vector_bbox_max_x],
		MaxY: b1[vector_bbox_max_y],
	}
}
