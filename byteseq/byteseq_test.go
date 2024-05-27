package byteseq

import (
	"testing"
)

func TestSettingBytesAsConsumed(t *testing.T) {
	// If we set 4 bytes when creating the random sequence, then
	// we should only be able to get 252 bytes out of the sequence
	// before it's exhausted
	testBytes := []byte{0x00, 0x01, 0x02, 0x03}
	expectedBytesCount := 256 - len(testBytes)

	byteSeq := NewRandomSeq(testBytes)
	bytesConsumed := 0

	for byteSeq.HasMore() {
		_, _ = byteSeq.NextValue()
		bytesConsumed++
	}

	if bytesConsumed != expectedBytesCount {
		t.Errorf("Expected only %d bytes but got %d", expectedBytesCount, bytesConsumed)
	}
}

// TODO: This can be improved using more values in a table-driven structure
func TestBasicMarkingOfConsumedValue(t *testing.T) {
	byteSeq := NewRandomSeq(nil)

	byteSeq.consumeByte(0x04)

	if byteSeq.consumedBytes[0] != 0x10 {
		t.Errorf("Expected byte 0x10 but got 0x%x", byteSeq.consumedBytes[0])
	}
}

func TestExhaustedSequenceReturnsError(t *testing.T) {
	byteSeq := NewRandomSeq(nil)

	for i := 0; i < 256; i++ {
		_, _ = byteSeq.NextValue()
	}

	// This one should return an error
	_, err := byteSeq.NextValue()
	if err == nil {
		t.Errorf("Expected error from exhausted byte sequence, but got nil")
	}
}

func TestNoByteValueIsRepeated(t *testing.T) {
	byteSeq := NewRandomSeq(nil)
	var returnedValueCounts [256]int

	for i := 0; i < 256; i++ {
		value, _ := byteSeq.NextValue()
		returnedValueCounts[value]++
	}

	// Now, every single entry should be one. Zero, or > 1 is an error
	for _, count := range returnedValueCounts {
		if count != 1 {
			t.Errorf("Expected only one occurrence of value but got %d", count)
		}
	}
}
