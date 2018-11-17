// func (n * Node32) vectorComputeDistances32 (x, y float32) (mind, maxd float32)
// +0 Minx
// +4 MinY
// +8 MaxX
// +12 MaxY
// +16 x
// +20 y
// +24 mind
// +28 maxd
TEXT ·vectorComputeDistances32(SB), $0-32
MOVSS  minx+0(FP), X0
MOVSS  miny+4(FP), X5
MOVSS  maxx+8(FP), X1
MOVSS  maxy+12(FP), X6
MOVSS pointx+16(FP), X2
MOVSS pointy+20(FP), X7

// compute for x
MOVSS X2, X3
SUBSS X0, X2 // point - min
SUBSS X1, X3 // point - max
MULSS X2, X2 // (point - min) ** 2
MULSS X3, X3 // (point - max) ** 2
MOVSS X2, X4 // copy to keep X2
MINSS X3, X4 // min of (point -min)**2, (point - max)**2
MAXSS X3, X2 // max of (point -min)**2, (point - max)**2
SUBSS X0, X1 // max - min
MULSS X1, X1 // (max - min)**2 (sides)
CMPSS X2, X1, 2// https://www.felixcloutier.com/x86/CMPPD.html sets 0 bits if false 1 bits if true order seems to be reversed
PAND X4, X1 // keep minx**2 if X1 is 1 mask (that is point is outside bbox)

// compute for y
MOVSS X7, X8
SUBSS X5, X7 // point - min
SUBSS X6, X8 // point - max
MULSS X7, X7 // (point - min) ** 2
MULSS X8, X8 // (point - max) ** 2
MOVSS X7, X9 // copy to keep X7
MINSS X8, X9 // min of (point -min)**2, (point - max)**2
MAXSS X8, X7 // max of (point -min)**2, (point - max)**2
SUBSS X5, X6 // max - min
MULSS X6, X6 // (max - min)**2 (sides)
CMPSS X7, X6, 2// https://www.felixcloutier.com/x86/CMPPD.html sets 0 bits if false 1 bits if true order seems to be reversed
PAND X9, X6 // keep minx if X1 is 1 mask

// MOVLPS X1, X3 // Move upper float of X1 to lower
ADDSS X6, X1 // Now X1 contains the sum of both elements
MOVSS X1, mind+24(FP)

// SHUFPD $0x1, X4, X4 // Invert X4
ADDSS X9, X2 // Crossed sum
ADDSS X4, X7
// MOVLPS X2, X6 // Move one of the mixed sums to X6

MINSS X2, X7 // Min of crossed sums
MOVSS X7, ret+28(FP)
RET
