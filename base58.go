// Package base58 implements an encoding scheme for binary data that represents
// the data using 58 alphanumeric characters.
package base58

import (
	"bytes"
	"math/big"
)

// All alphanumeric characters except for "0", "I", "O", and "l"
var alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Encode(src []byte) string {
	nZeroes := bytes.IndexFunc(src, func(r rune) bool { return r != rune(0) })
	if nZeroes != -1 {
		src = src[nZeroes:]
	} else {
		nZeroes = len(src)
	}

	// Allocate enough capacity for worst case scenario, e.g. log(256) / log(58), rounded up.
	sz := len(src)*137/100 + nZeroes + 1
	output := make([]rune, sz)

	for i := 0; i < nZeroes; i++ {
		output[i] = '1'
	}

	srcInt := new(big.Int).SetBytes(src)
	targetInt := new(big.Int)
	// Start at the end, work our way back, since we're going to produce the encoded string in reverse.
	newIdx := sz - 1
	for srcInt.Sign() > 0 {
		srcInt.DivMod(srcInt, big.NewInt(58), targetInt)
		output[newIdx] = rune(alphabet[targetInt.Int64()])
		newIdx--
	}

	return string(append(output[:nZeroes], output[newIdx+1:]...))
}

func Decode(src string) ([]byte, bool) {
	dat := []byte(src)
	nZeroes := bytes.IndexFunc(dat, func(r rune) bool { return r != '1' })
	if nZeroes != -1 {
		src = src[nZeroes:]
	} else {
		nZeroes = len(src)
	}
	srcInt := new(big.Int)
	tmpbig := new(big.Int)
	for _, v := range dat {
		idx := bytes.IndexByte(alphabet, v)
		if idx == -1 {
			return nil, false
		}
		tmpbig.SetInt64(int64(idx))
		srcInt.Mul(srcInt, big.NewInt(58))
		srcInt.Add(srcInt, tmpbig)
	}
	targetBytes := srcInt.Bytes()
	if nZeroes > 0 {
		targetBytes = append(make([]byte, nZeroes), targetBytes...)
	}
	return targetBytes, true
}
