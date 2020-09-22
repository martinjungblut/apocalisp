package core

import (
	"testing"
)

func Test_NewHashmap_IsHashmap_IsEmptyHashmap(t *testing.T) {
	hashmap := NewHashmap()

	if !hashmap.IsHashmap() {
		t.Error("NewHashmap() failed.")
	}

	if !hashmap.IsEmptyHashmap() {
		t.Error("NewHashmap() failed.")
	}
}

func Test_NewHashmapFromSequence_With_Even_Values(t *testing.T) {
	first, second := *NewString("first"), *NewString("second")
	third, fourth := *NewString("third"), *NewString("fourth")

	sequence := make([]Type, 0)
	sequence = append(sequence, first)
	sequence = append(sequence, second)
	sequence = append(sequence, third)
	sequence = append(sequence, fourth)

	hashmap := NewHashmapFromSequence(sequence).AsHashmap()

	if hashmap[first] != second {
		t.Error("NewHashmapFromSequence() failed.")
	}

	if hashmap[third] != fourth {
		t.Error("NewHashmapFromSequence() failed.")
	}

	if len(hashmap) != 2 {
		t.Errorf("NewHashmapFromSequence() failed. Length should be 2, but it's '%d'.", len(hashmap))
	}
}

func Test_NewHashmapFromSequence_With_Odd_Values(t *testing.T) {
	first, second, third := *NewString("first"), *NewString("second"), *NewString("third")

	sequence := make([]Type, 0)
	sequence = append(sequence, first)
	sequence = append(sequence, second)
	sequence = append(sequence, third)

	hashmap := NewHashmapFromSequence(sequence).AsHashmap()

	if hashmap[first] != second {
		t.Error("NewHashmapFromSequence() failed.")
	}

	if len(hashmap) != 1 {
		t.Errorf("NewHashmapFromSequence() failed. Length should be 1, but it's '%d'.", len(hashmap))
	}
}
