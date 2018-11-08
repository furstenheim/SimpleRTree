// func (n * Node) vectorComputeDistances (x, y float64) (mind, maxd float64)
// +0 Minx
// +8 MinY
// +16 MaxX
// +24 MaxY
// +32 x
// +40 y
// +48 mind
// +56 maxd
TEXT ·vectorComputeDistances(SB), $0-60
MOVSD  minx+0(FP), X0
MOVSD  miny+8(FP), X5
MOVSD  maxx+16(FP), X1
MOVSD  maxy+24(FP), X6
MOVSD pointx+32(FP), X2
MOVSD pointy+40(FP), X7

// compute for x
MOVSD X2, X3
SUBSD X0, X2 // point - min
SUBSD X1, X3 // point - max
MULSD X2, X2 // (point - min) ** 2
MULSD X3, X3 // (point - max) ** 2
MOVSD X2, X4 // copy to keep X2
MINSD X3, X4 // min of (point -min)**2, (point - max)**2
MAXSD X3, X2 // max of (point -min)**2, (point - max)**2
SUBSD X0, X1 // max - min
MULSD X1, X1 // (max - min)**2 (sides)
CMPSD X2, X1, 2// https://www.felixcloutier.com/x86/CMPPD.html sets 0 bits if false 1 bits if true order seems to be reversed
PAND X4, X1 // keep minx if X1 is 1 mask

// compute for y
MOVSD X7, X8
SUBSD X5, X7 // point - min
SUBSD X6, X8 // point - max
MULSD X7, X7 // (point - min) ** 2
MULSD X8, X8 // (point - max) ** 2
MOVSD X7, X9 // copy to keep X7
MINSD X8, X9 // min of (point -min)**2, (point - max)**2
MAXSD X8, X7 // max of (point -min)**2, (point - max)**2
SUBSD X5, X6 // max - min
MULSD X6, X6 // (max - min)**2 (sides)
CMPSD X7, X6, 2// https://www.felixcloutier.com/x86/CMPPD.html sets 0 bits if false 1 bits if true order seems to be reversed
PAND X9, X6 // keep minx if X1 is 1 mask

// MOVLPS X1, X3 // Move upper float of X1 to lower
ADDSD X6, X1 // Now X1 contains the sum of both elements
MOVSD X1, mind+48(FP)

SHUFPD $0x1, X4, X4 // Invert X4
ADDSD X4, X2 // Crossed sum
ADDSD X9, X7
// MOVLPS X2, X6 // Move one of the mixed sums to X6

MINSD X2, X7 // Min of crossed sums
MOVSD X7, ret+56(FP)
RET
