package main

import (
	"strings"
	"testing"
)

func TestFindAnagrams(t *testing.T) {
	data := []string{"Пятак", "пятка", "пятка", "тяпка", "листок", "слиток", "столик", "КЛОУН", "Колун", "уклон", "кулон"}
	expected := map[string][]string{
		"пятак":  {"пятак", "пятка", "тяпка"},
		"листок": {"листок", "слиток", "столик"},
		"клоун":  {"клоун", "колун", "кулон", "уклон"},
	}
	anagrams := FindAnagrams(data)
	for key, val := range anagrams {
		anagram, ok := expected[key]
		if !ok {
			t.Errorf("key %s not found", key)
		}
		a := strings.Join(anagram, " ")
		b := strings.Join(val, " ")
		if a != b {
			t.Errorf("\nIncorrect anagrams:\nResult:\n  %s\nExpected:\n  %s", b, a)
		}
	}
}
