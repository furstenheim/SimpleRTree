// func sse2Vec64Add(a, b Vec64) Vec64
// a: +0(FP)
// b: +32(FP)
// Return: +64(FP)
TEXT ·sse2Vec64Add(SB),$0-128
	MOVUPD    a+0(FP), X0
	MOVUPD    b+32(FP), X2
	ADDPD     X2, X0
	MOVUPD    X0, ret+64(FP)
	MOVUPD    a+16(FP), X0
	MOVUPD    b+48(FP), X2
	ADDPD     X2, X0
	MOVUPD    X0, ret+80(FP)
	RET


// func sse2Vec64Sub(a, b Vec64) Vec64
// a: +0(FP)
// b: +32(FP)
// Return: +64(FP)
TEXT ·sse2Vec64Sub(SB),$0-128
	MOVUPD    a+0(FP), X0
	MOVUPD    b+32(FP), X2
	SUBPD     X2, X0
	MOVUPD    X0, ret+64(FP)
	MOVUPD    a+16(FP), X0
	MOVUPD    b+48(FP), X2
	SUBPD     X2, X0
	MOVUPD    X0, ret+80(FP)
	RET

// func avxVec64Sub(a, b Vec64) Vec64
// a: +0(SP)
// b: +32(SP)
// Return: +64(FP)
TEXT ·avxVec64Sub(SB),$0-128
	BYTE $0xc5; BYTE $0xfd; BYTE $0x10; BYTE $0x44; BYTE $0x24; BYTE $0x08; // vmovupd 0x8(%rsp),%ymm0
	BYTE $0xc5; BYTE $0xfd; BYTE $0x5c; BYTE $0x4c; BYTE $0x24; BYTE $0x28; // vsubpd 0x28(%rsp),%ymm0,%ymm1
	BYTE $0xc5; BYTE $0xfd; BYTE $0x11; BYTE $0x4c; BYTE $0x24; BYTE $0x48; // vmovupd %ymm1,0x48(%rsp)
	RET


// func sse2Vec64Mul(a, b Vec64) Vec64
// a: +0(FP)
// b: +32(FP)
// Return: +64(FP)
TEXT ·sse2Vec64Mul(SB),$0-128
	MOVUPD    a+0(FP), X0
	MOVUPD    b+32(FP), X2
	MULPD     X2, X0
	MOVUPD    X0, ret+64(FP)
	MOVUPD    a+16(FP), X0
	MOVUPD    b+48(FP), X2
	MULPD     X2, X0
	MOVUPD    X0, ret+80(FP)
	RET

// func avxVec64Mul(a, b Vec64) Vec64
// a: +0(SP)
// b: +32(SP)
// Return: +64(FP)
TEXT ·avxVec64Mul(SB),$0-128
	BYTE $0xc5; BYTE $0xfd; BYTE $0x10; BYTE $0x44; BYTE $0x24; BYTE $0x08; // vmovupd 0x8(%rsp),%ymm0
	BYTE $0xc5; BYTE $0xfd; BYTE $0x59; BYTE $0x4c; BYTE $0x24; BYTE $0x28; // vmulpd 0x28(%rsp),%ymm0,%ymm1
	BYTE $0xc5; BYTE $0xfd; BYTE $0x11; BYTE $0x4c; BYTE $0x24; BYTE $0x48; // vmovupd %ymm1,0x48(%rsp)
	RET
