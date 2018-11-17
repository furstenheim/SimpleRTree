package SimpleRTree

type VectorBBox32 [4]float32


func newVectorBBox32(MinX, MinY, MaxX, MaxY float32) VectorBBox32 {
	return [4]float32{MinX, MinY, MaxX, MaxY}
}

/*
func bbox2VectorBBox32(b BBox) VectorBBox32 {
	return newVectorBBox32(b.MinX, b.MinY, b.MaxX, b.MaxY)
}
*/

func vectorBBoxExtend32(b1, b2 VectorBBox32) VectorBBox32{
	return [4]float32{
		minFloat32(b1[0], b2[0]),
		minFloat32(b1[1], b2[1]),
		maxFloat32(b1[2], b2[2]),
		maxFloat32(b1[3], b2[3]),
	}
}
/*

func (b1 VectorBBox32) toBBox() BBox {
	return BBox{
		MinX: b1[VECTOR_BBOX_MIN_X],
		MinY: b1[VECTOR_BBOX_MIN_Y],
		MaxX: b1[VECTOR_BBOX_MAX_X],
		MaxY: b1[VECTOR_BBOX_MAX_Y],
	}
}
*/
