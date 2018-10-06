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
MOVUPD point+32(FP), X2
MOVUPD  min+0(FP), X0
MOVUPD  max+16(FP), X1
MOVAPS X2, X3
SUBPD X0, X2 // point - min
SUBPD X1, X3 // point - max
MULPD X2, X2 // (point - min) ** 2
MULPD X3, X3 // (point - max) ** 2
MOVAPS X2, X4 // copy to keep X2
MINPD X3, X4 // min of (point -min)**2, (point - max)**2
MAXPD X3, X2 // max of (point -min)**2, (point - max)**2
SUBPD X0, X1 // max - min
MULPD X1, X1 // (max - min)**2 (sides)
CMPPD X2, X1, 2// https://www.felixcloutier.com/x86/CMPPD.html sets 0 bits if false 1 bits if true order seems to be reversed
ANDPD X4, X1 // keep minx if X1 is 1 mask

MOVLPS X1, X3 // Move upper float of X1 to lower
ADDSD X3, X1 // Now X1 contains the sum of both elements
MOVLPS X1, mind+48(FP)


MOVLPS X4, X3
ADDSD X2, X3 // X3 contains min1 + max2

MOVLPS X2, X5
ADDSD X4, X5

MINSD X3, X5
MOVLPS X5, ret+56(FP)
RET
