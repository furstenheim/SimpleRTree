package SimpleRTree

import "math"

// Original code from https://mmcloughlin.com/posts/geohash-assembly
func GeoHash(lat, lng float64) uint64 {
	return interleave(hashQuantize(lat, lng))
}

func hashQuantize(lat, lng float64) (lat32 uint32, lng32 uint32) {
	lat32 = uint32(math.Ldexp((lat+90.0)/180.0, 32))
	lng32 = uint32(math.Ldexp((lng+180.0)/360.0, 32))
	return
}

func spread(x uint32) uint64 {
	X := uint64(x)
	X = (X | (X << 16)) & 0x0000ffff0000ffff
	X = (X | (X << 8)) & 0x00ff00ff00ff00ff
	X = (X | (X << 4)) & 0x0f0f0f0f0f0f0f0f
	X = (X | (X << 2)) & 0x3333333333333333
	X = (X | (X << 1)) & 0x5555555555555555
	return X
}

func interleave(x, y uint32) uint64 {
	return spread(x) | (spread(y) << 1)
}

type GeoHashSorter struct {
	points FlatPoints
	hashes []uint64
}


func (s GeoHashSorter) Less(i, j int) bool {
	return s.hashes[i] < s.hashes[j]
}

func (s GeoHashSorter) Swap(i, j int) {
	s.points.Swap(i, j)
	s.hashes[i], s.hashes[j] = s.hashes[j], s.hashes[i]
}

func (s GeoHashSorter) Len() int {
	return len(s.hashes)
}
