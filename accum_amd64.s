//go:build amd64

#include "textflag.h"

// Secret constants (must match rapidhash.go)
#define SECRET0 $0x2d358dccaa6c78a5
#define SECRET1 $0x8bb84b93962eacc9
#define SECRET2 $0x4b33a62ed433d4a3
#define SECRET3 $0x4d5a2da51de1aa47
#define SECRET4 $0xa0761d6478bd642f
#define SECRET5 $0xe7037ed1a0b428db
#define SECRET6 $0x90ed1765281c388c

// func accumBlocks(p unsafe.Pointer, length int, seed, see1, see2, see3, see4, see5, see6 uint64) (
//     newP unsafe.Pointer, remaining int, nseed, nsee1, nsee2, nsee3, nsee4, nsee5, nsee6 uint64)
//
// Processes multiple 112-byte blocks with loop unrolling (2 blocks = 224 bytes per iteration).
// This version keeps all accumulators in registers throughout the loop.
//
TEXT ·accumBlocks(SB), NOSPLIT, $0-144
    // Load inputs
    MOVQ p+0(FP), SI          // SI = data pointer
    MOVQ length+8(FP), CX     // CX = remaining length
    
    // Load accumulators into callee-saved registers
    MOVQ seed+16(FP), DI
    MOVQ see1+24(FP), R8
    MOVQ see2+32(FP), R9
    MOVQ see3+40(FP), R10
    MOVQ see4+48(FP), R11
    MOVQ see5+56(FP), R12
    MOVQ see6+64(FP), R13

    // Check if we have at least 224 bytes for unrolled loop
    CMPQ CX, $224
    JLE single_block_loop

unrolled_loop:    
    // mix 0: seed = mix(p[0]^secret0, p[8]^seed)
    MOVQ 0(SI), AX
    MOVQ 8(SI), BX
    MOVQ SECRET0, R14
    XORQ R14, AX
    XORQ DI, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, DI

    // mix 1: see1 = mix(p[16]^secret1, p[24]^see1)
    MOVQ 16(SI), AX
    MOVQ 24(SI), BX
    MOVQ SECRET1, R14
    XORQ R14, AX
    XORQ R8, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R8

    // mix 2: see2 = mix(p[32]^secret2, p[40]^see2)
    MOVQ 32(SI), AX
    MOVQ 40(SI), BX
    MOVQ SECRET2, R14
    XORQ R14, AX
    XORQ R9, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R9

    // mix 3: see3 = mix(p[48]^secret3, p[56]^see3)
    MOVQ 48(SI), AX
    MOVQ 56(SI), BX
    MOVQ SECRET3, R14
    XORQ R14, AX
    XORQ R10, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R10

    // mix 4: see4 = mix(p[64]^secret4, p[72]^see4)
    MOVQ 64(SI), AX
    MOVQ 72(SI), BX
    MOVQ SECRET4, R14
    XORQ R14, AX
    XORQ R11, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R11

    // mix 5: see5 = mix(p[80]^secret5, p[88]^see5)
    MOVQ 80(SI), AX
    MOVQ 88(SI), BX
    MOVQ SECRET5, R14
    XORQ R14, AX
    XORQ R12, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R12

    // mix 6: see6 = mix(p[96]^secret6, p[104]^see6)
    MOVQ 96(SI), AX
    MOVQ 104(SI), BX
    MOVQ SECRET6, R14
    XORQ R14, AX
    XORQ R13, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R13

    
    // mix 0: seed = mix(p[112]^secret0, p[120]^seed)
    MOVQ 112(SI), AX
    MOVQ 120(SI), BX
    MOVQ SECRET0, R14
    XORQ R14, AX
    XORQ DI, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, DI

    // mix 1: see1 = mix(p[128]^secret1, p[136]^see1)
    MOVQ 128(SI), AX
    MOVQ 136(SI), BX
    MOVQ SECRET1, R14
    XORQ R14, AX
    XORQ R8, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R8

    // mix 2: see2 = mix(p[144]^secret2, p[152]^see2)
    MOVQ 144(SI), AX
    MOVQ 152(SI), BX
    MOVQ SECRET2, R14
    XORQ R14, AX
    XORQ R9, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R9

    // mix 3: see3 = mix(p[160]^secret3, p[168]^see3)
    MOVQ 160(SI), AX
    MOVQ 168(SI), BX
    MOVQ SECRET3, R14
    XORQ R14, AX
    XORQ R10, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R10

    // mix 4: see4 = mix(p[176]^secret4, p[184]^see4)
    MOVQ 176(SI), AX
    MOVQ 184(SI), BX
    MOVQ SECRET4, R14
    XORQ R14, AX
    XORQ R11, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R11

    // mix 5: see5 = mix(p[192]^secret5, p[200]^see5)
    MOVQ 192(SI), AX
    MOVQ 200(SI), BX
    MOVQ SECRET5, R14
    XORQ R14, AX
    XORQ R12, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R12

    // mix 6: see6 = mix(p[208]^secret6, p[216]^see6)
    MOVQ 208(SI), AX
    MOVQ 216(SI), BX
    MOVQ SECRET6, R14
    XORQ R14, AX
    XORQ R13, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R13

    // Advance pointer and decrement length
    ADDQ $224, SI
    SUBQ $224, CX
    
    // Continue if we have 224+ bytes remaining
    CMPQ CX, $224
    JG unrolled_loop

single_block_loop:
    // Process remaining 112-byte blocks one at a time
    CMPQ CX, $112
    JLE done

    // mix 0
    MOVQ 0(SI), AX
    MOVQ 8(SI), BX
    MOVQ SECRET0, R14
    XORQ R14, AX
    XORQ DI, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, DI

    // mix 1
    MOVQ 16(SI), AX
    MOVQ 24(SI), BX
    MOVQ SECRET1, R14
    XORQ R14, AX
    XORQ R8, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R8

    // mix 2
    MOVQ 32(SI), AX
    MOVQ 40(SI), BX
    MOVQ SECRET2, R14
    XORQ R14, AX
    XORQ R9, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R9

    // mix 3
    MOVQ 48(SI), AX
    MOVQ 56(SI), BX
    MOVQ SECRET3, R14
    XORQ R14, AX
    XORQ R10, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R10

    // mix 4
    MOVQ 64(SI), AX
    MOVQ 72(SI), BX
    MOVQ SECRET4, R14
    XORQ R14, AX
    XORQ R11, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R11

    // mix 5
    MOVQ 80(SI), AX
    MOVQ 88(SI), BX
    MOVQ SECRET5, R14
    XORQ R14, AX
    XORQ R12, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R12

    // mix 6
    MOVQ 96(SI), AX
    MOVQ 104(SI), BX
    MOVQ SECRET6, R14
    XORQ R14, AX
    XORQ R13, BX
    MULQ BX
    XORQ DX, AX
    MOVQ AX, R13

    ADDQ $112, SI
    SUBQ $112, CX
    JMP single_block_loop

done:
    // Store results
    MOVQ SI, newP+72(FP)
    MOVQ CX, remaining+80(FP)
    MOVQ DI, nseed+88(FP)
    MOVQ R8, nsee1+96(FP)
    MOVQ R9, nsee2+104(FP)
    MOVQ R10, nsee3+112(FP)
    MOVQ R11, nsee4+120(FP)
    MOVQ R12, nsee5+128(FP)
    MOVQ R13, nsee6+136(FP)
    RET
