package byteseq

import (
	"fmt"
	"math/rand"
)

type RandomByteSeq struct {
	consumedBytes   [32]byte // Bitmap for each of 256 possible values
	remainingValues int
}

func NewRandomSeq(consumedBytes []byte) *RandomByteSeq {
	seq := &RandomByteSeq{}

	seq.remainingValues = 256

	// Were any already-consumed bytes specified?
	if len(consumedBytes) > 0 {
		for _, b := range consumedBytes {
			seq.consumeByte(b)
		}
	}
	return seq
}

func (r *RandomByteSeq) valueHasBeenConsumed(b byte) bool {
	// Get the correct byte of our bitmap from top 5 bits
	byteIndex := b & 0xF8 >> 3
	testByte := r.consumedBytes[byteIndex]

	// Now, from the bottom 3 bits, check the appropriate bit
	bitMask := byte(1 << (b & 0x07))

	return testByte&bitMask != 0
}

func (r *RandomByteSeq) consumeByte(b byte) {
	byteIndex := b & 0xF8 >> 3
	bitMask := byte(1 << (b & 0x07))

	r.consumedBytes[byteIndex] |= bitMask
	r.remainingValues--
}

func (r *RandomByteSeq) HasMore() bool {
	return r.remainingValues > 0
}

func (r *RandomByteSeq) NextValue() (byte, error) {
	// Are there any more values available?
	if r.remainingValues == 0 {
		return 0, fmt.Errorf("sequence has been exhausted")
	}

	for {
		valueByte := byte(rand.Int() & 0xFF)

		if !r.valueHasBeenConsumed(valueByte) {
			r.consumeByte(valueByte)
			return valueByte, nil
		}
	}
}
