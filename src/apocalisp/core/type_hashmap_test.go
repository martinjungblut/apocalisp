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

	if hashmap["first"] != second {
		t.Error("NewHashmapFromSequence() failed.")
	}

	if hashmap["third"] != fourth {
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

	if hashmap["first"] != second {
		t.Error("NewHashmapFromSequence() failed.")
	}

	if len(hashmap) != 1 {
		t.Errorf("NewHashmapFromSequence() failed. Length should be 1, but it's '%d'.", len(hashmap))
	}
}

func Test_NewHashmapFromSequence_With_Mixed_Strings_And_Keywords(t *testing.T) {
	first, second, third, fourth, fifth, sixth := *NewSymbol(":first"), *NewBoolean(true), *NewString(":third"), *NewBoolean(false), *NewSymbol("fifth"), *NewBoolean(true)

	sequence := make([]Type, 0)
	sequence = append(sequence, first)
	sequence = append(sequence, second)
	sequence = append(sequence, third)
	sequence = append(sequence, fourth)
	sequence = append(sequence, fifth)
	sequence = append(sequence, sixth)

	hashmap := NewHashmapFromSequence(sequence).AsHashmap()

	hfirst, hthird, hfifth := hashmap[":first"], hashmap[":third"], hashmap["fifth"]
	if !hfirst.AsBoolean() {
		t.Error("NewHashmapFromSequence() failed.")
	}
	if !hfirst.HashmapSymbolValue || second.HashmapSymbolValue {
		t.Error("NewHashmapFromSequence() failed.")
	}

	if hthird.AsBoolean() {
		t.Error("NewHashmapFromSequence() failed.")
	}
	if hthird.HashmapSymbolValue || fourth.HashmapSymbolValue {
		t.Error("NewHashmapFromSequence() failed.")
	}

	if !hfifth.AsBoolean() {
		t.Error("NewHashmapFromSequence() failed.")
	}
	if !hfifth.HashmapSymbolValue || sixth.HashmapSymbolValue {
		t.Error("NewHashmapFromSequence() failed.")
	}
}
