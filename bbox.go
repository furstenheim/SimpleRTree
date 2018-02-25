package SimpleRTree

import (
	"math"
)

type BBox struct {
	MinX, MinY, MaxX, MaxY float64
}

func (b BBox) area() float64 {
	return (b.MaxX - b.MinX) * (b.MaxY - b.MinY)
}

func (b1 BBox) equals (b2 BBox) bool {
	return b1.MinX == b2.MinX &&
		b1.MinY == b2.MinY &&
		b1.MaxX == b2.MaxX &&
		b1.MaxY == b2.MaxY
}
func (b1 BBox) extend(b2 BBox) BBox {
	return BBox{
		MinX: math.Min(b1.MinX, b2.MinX),
		MinY: math.Min(b1.MinY, b2.MinY),
		MaxX: math.Max(b1.MaxX, b2.MaxX),
		MaxY: math.Max(b1.MaxY, b2.MaxY),
	}
}

func (b1 BBox) intersectionArea(b2 BBox) float64 {
	minX := math.Max(b1.MinX, b2.MinX)
	maxX := math.Min(b1.MaxX, b2.MaxX)
	minY := math.Max(b1.MinY, b2.MinY)
	maxY := math.Min(b1.MaxY, b2.MaxY)
	return math.Max(0, maxX-minX) * math.Min(0, maxY-minY)
}

func (b1 BBox) contains(b2 BBox) bool {
	return b1.MinX <= b2.MinX &&
		b2.MaxX <= b1.MaxX &&
		b1.MinY <= b2.MinY &&
		b2.MaxY <= b1.MaxY
}

func (b1 BBox) containsPoint (x, y float64) bool {
	return b1.MinX <= x &&
		x <= b1.MaxX &&
		b1.MinY <= y &&
		y <= b1.MaxY
}



func (b1 BBox) intersects(b2 BBox) bool {
	return b2.MinX <= b1.MaxX &&
		b2.MinY <= b1.MaxY &&
		b2.MaxX >= b1.MinX &&
		b2.MaxY >= b1.MinY
}

func (b1 BBox) enlargedArea(b2 BBox) float64 {
	return (math.Max(b2.MaxX, b1.MaxX) - math.Min(b2.MinX, b1.MinX)) *
		(math.Max(b2.MaxY, b1.MaxY) - math.Min(b2.MinY, b1.MinY))
}
