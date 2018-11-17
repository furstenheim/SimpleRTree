// func (n * Node) vectorBBoxExtendASM (b1, b2 VectorBBox) (VectorBBox)
// +0 b1.x1
// +8 b1.y1
// +16 b1.x2
// +24 b1.y2
// +32 b2.x1
// +40 b2.y1
// +48 b2.x2
// +56 b2.y2
// +64
// +72
// +80
// +88
TEXT ·vectorBBoxExtendASM(SB), $0-60
MOVSD  b1x1+0(FP), X0
MOVSD  b1y1+8(FP), X1
MOVSD  b1y2+16(FP), X2
MOVSD  b1x2+24(FP), X3
MINSD  b2x1+32(FP), X0
MINSD  b2y1+40(FP), X1
MAXSD  b2y2+48(FP), X2
MAXSD  b2x2+56(FP), X3
MOVSD  X0, ret+64(FP)
MOVSD  X1, ret+72(FP)
MOVSD  X2, ret+80(FP)
MOVSD  X3, ret+88(FP)
RET
