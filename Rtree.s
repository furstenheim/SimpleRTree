// func (n * Node) vectorComputeDistances (x, y float64) (mind, maxd float64)
// +0 Minx
// +8 MinY
// +16 MaxX
// +24 MaxY
// +32 x
// +40 y
// +48 mind
// +54 maxd
// Return +24, +32
TEXT ·vectorComputeDistances(SB), $0-60
MOVUPD  min+0(FP), X0
MOVUPD  max+16(FP), X1
MOVUPD point+32(FP), X2
MOVUPD point+32(FP), X3
SUBPD X0, X2 // point - min
SUBPD X1, X3 // point - max
MULPD X2, X2 // (point - min) ** 2
MULPD X3, X3 // (point - max) ** 2
MOVUPD X2, X4 // copy to keep X2
MOVUPD X2, X5 // copy to keep X2
MINPD X3, X4 // min of (point -min)**2, (point - max)**2
MAXPD X3, X5 // max of (point -min)**2, (point - max)**2
SUBPD X0, X1 // max - min
MULPD X1, X1 // (max - min)**2

MOVUPD X1, ret+48(FP)
RET

    /*
        0x0000 00000 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    MOVQ    "".n+8(SP), AX
        0x0005 00005 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    MOVSD   16(AX), X0
        0x000a 00010 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    MOVSD   "".x+16(SP), X1
        0x0010 00016 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    MOVUPS  X1, X2
        0x0013 00019 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    SUBSD   X0, X1
        0x0017 00023 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    MULSD   X1, X1
        0x001b 00027 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    MOVSD   32(AX), X3
        0x0020 00032 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    SUBSD   X3, X2
        0x0024 00036 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    MULSD   X2, X2
        0x0028 00040 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    UCOMISD X2, X1
        0x002c 00044 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    JLS     179
        0x0032 00050 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    MOVSD   24(AX), X4
        0x0037 00055 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    MOVSD   "".y+24(SP), X5
        0x003d 00061 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    MOVUPS  X5, X6
        0x0040 00064 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    SUBSD   X4, X5
        0x0044 00068 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    MULSD   X5, X5
        0x0048 00072 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    PCDATA  $2, $0
        0x0048 00072 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    MOVSD   40(AX), X7
        0x004d 00077 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    SUBSD   X7, X6
        0x0051 00081 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    MULSD   X6, X6
        0x0055 00085 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    UCOMISD X6, X5
        0x0059 00089 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    JLS     166
        0x005b 00091 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:421)    SUBSD   X0, X3
        0x005f 00095 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:421)    MULSD   X3, X3
        0x0063 00099 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:422)    SUBSD   X4, X7
        0x0067 00103 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:422)    MULSD   X7, X7
        0x006b 00107 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:425)    UCOMISD X3, X1
        0x006f 00111 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:425)    JCS     158
        0x0071 00113 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    MOVUPS  X2, X0
        0x0074 00116 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:428)    UCOMISD X7, X5
        0x0078 00120 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:428)    JCS     126
        0x007a 00122 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:429)    ADDSD   X6, X2
        0x007e 00126 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:433)    ADDSD   X1, X6
        0x0082 00130 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:433)    ADDSD   X0, X5
        0x0086 00134 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:433)    UCOMISD X6, X5
        0x008a 00138 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:433)    JLS     153
        0x008c 00140 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:434)    MOVSD   X2, "".mind+32(SP)
        0x0092 00146 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:434)    MOVSD   X6, "".maxd+40(SP)
        0x0098 00152 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:434)    RET
        0x0099 00153 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:433)    MOVUPS  X5, X6
        0x009c 00156 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:433)    JMP     140
        0x009e 00158 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:433)    MOVUPS  X2, X0
        0x00a1 00161 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:433)    XORPS   X2, X2
        0x00a4 00164 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:425)    JMP     116
        0x00a6 00166 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    MOVUPS  X6, X8
        0x00aa 00170 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    MOVUPS  X5, X6
        0x00ad 00173 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    MOVUPS  X8, X5
        0x00b1 00177 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:419)    JMP     91
        0x00b3 00179 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    PCDATA  $2, $1
        0x00b3 00179 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    MOVUPS  X2, X4
        0x00b6 00182 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    MOVUPS  X1, X2
        0x00b9 00185 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    MOVUPS  X4, X1
        0x00bc 00188 (/home/gabi/Gabi/apps/SimpleRTree/RTree.go:418)    JMP     50
*/