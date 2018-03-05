package SimpleRTree

// Code from https://github.com/slimsag/rand/blob/master/simd/vec64.go since original repo is deprecated
type vec64 [4]float64

// Implemented in vec64.s
func sse2Vec64Add(a, b vec64) vec64

// Implemented in vec64.s
func sse2Vec64Sub(a, b vec64) vec64

// Implemented in vec64.s
func avxVec64Sub(a, b vec64) vec64

// Implemented in vec64.s
func sse2Vec64Mul(a, b vec64) vec64


// Implemented in vec64.s
func avxVec64Mul(a, b vec64) vec64