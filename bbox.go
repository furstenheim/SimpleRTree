package SimpleRTree

import (
	"math"
)

type rBBox struct {
	MinX, MinY, MaxX, MaxY float64
}

func (b rBBox) area() float64 {
	return (b.MaxX - b.MinX) * (b.MaxY - b.MinY)
}

func (b1 rBBox) equals(b2 rBBox) bool {
	return b1.MinX == b2.MinX &&
		b1.MinY == b2.MinY &&
		b1.MaxX == b2.MaxX &&
		b1.MaxY == b2.MaxY
}
func (b1 rBBox) extend(b2 rBBox) rBBox {
	return rBBox{
		MinX: math.Min(b1.MinX, b2.MinX),
		MinY: math.Min(b1.MinY, b2.MinY),
		MaxX: math.Max(b1.MaxX, b2.MaxX),
		MaxY: math.Max(b1.MaxY, b2.MaxY),
	}
}

func (b1 rBBox) intersectionArea(b2 rBBox) float64 {
	minX := math.Max(b1.MinX, b2.MinX)
	maxX := math.Min(b1.MaxX, b2.MaxX)
	minY := math.Max(b1.MinY, b2.MinY)
	maxY := math.Min(b1.MaxY, b2.MaxY)
	return math.Max(0, maxX-minX) * math.Min(0, maxY-minY)
}

func (b1 rBBox) contains(b2 rBBox) bool {
	return b1.MinX <= b2.MinX &&
		b2.MaxX <= b1.MaxX &&
		b1.MinY <= b2.MinY &&
		b2.MaxY <= b1.MaxY
}

func (b1 rBBox) containsPoint(x, y float64) bool {
	return b1.MinX <= x &&
		x <= b1.MaxX &&
		b1.MinY <= y &&
		y <= b1.MaxY
}

func (b1 rBBox) intersects(b2 rBBox) bool {
	return b2.MinX <= b1.MaxX &&
		b2.MinY <= b1.MaxY &&
		b2.MaxX >= b1.MinX &&
		b2.MaxY >= b1.MinY
}

func (b1 rBBox) enlargedArea(b2 rBBox) float64 {
	return (math.Max(b2.MaxX, b1.MaxX) - math.Min(b2.MinX, b1.MinX)) *
		(math.Max(b2.MaxY, b1.MaxY) - math.Min(b2.MinY, b1.MinY))
}
