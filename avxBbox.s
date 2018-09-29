// func avxBBoxExtendAux(b1, b2 AvxBBox) AvxBBox
// b1: +0(SP)
// b2: +32(SP)
// Return: +64(FP)
TEXT ·avxBBoxExtendAvx(SB),$0-128
BYTE $0xc5; BYTE $0xfd; BYTE $0x10; BYTE $0x44; BYTE $0x24; BYTE $0x08; // vmovupd 0x8(%rsp),%ymm0
BYTE $0xc5; BYTE $0xfd; BYTE $0x5d; BYTE $0x4c; BYTE $0x24; BYTE $0x28; // vminpd 0x28(%rsp),%ymm0,%ymm1
BYTE $0xc5; BYTE $0xfd; BYTE $0x11; BYTE $0x4c; BYTE $0x24; BYTE $0x48; // vmovupd %ymm1,0x48(%rsp)
RET

// func avxBBoxExtend(b1, b2 AvxBBox) AvxBBox
// b1: +0(SP)
// b2: +32(SP)
// Return: +64(FP)
TEXT ·avxBBoxExtend(SB),$0-128
MOVUPD    a+0(FP), X0
MOVUPD    b+32(FP), X2
MINPD     X2, X0
MOVUPD    X0, ret+64(FP)
MOVUPD    a+16(FP), X0
MOVUPD    b+48(FP), X2
MINPD     X2, X0
MOVUPD    X0, ret+80(FP)
RET
