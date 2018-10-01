// func vectorBBoxExtend(b1, b2 VectorBBox) VectorBBox
// b1: +0(SP)
// b2: +32(SP)
// Return: +64(FP)
TEXT ·vectorBBoxExtend(SB),$0-128
MOVUPD    a+0(FP), X0
MOVUPD    b+32(FP), X2
MINPD     X2, X0
MOVUPD    X0, ret+64(FP)
MOVUPD    a+16(FP), X0
MOVUPD    b+48(FP), X2
MAXPD     X2, X0
MOVUPD    X0, ret+80(FP)
RET
