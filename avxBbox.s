// func avxBBoxExtend(b1, b2 AvxBBox) AvxBBox
// b1: +0(SP)
// b2: +32(SP)
// Return: +64(FP)
TEXT ·avxBBoxExtend(SB),$0-128
BYTE $0xc5; BYTE $0xfd; BYTE $0x10; BYTE $0x44; BYTE $0x24; BYTE $0x08; // vmovupd 0x8(%rsp),%ymm0
BYTE $0xc5; BYTE $0xfd; BYTE $0x5d; BYTE $0x4c; BYTE $0x24; BYTE $0x28; // vminpd 0x28(%rsp),%ymm0,%ymm1
BYTE $0xc5; BYTE $0xfd; BYTE $0x11; BYTE $0x4c; BYTE $0x24; BYTE $0x48; // vmovupd %ymm1,0x48(%rsp)
RET
