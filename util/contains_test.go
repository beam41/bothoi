package util

import "testing"

func Test_ContainsStr_ReturnTrueWhenStringFoundInSlice(t *testing.T) {
	result := ContainsStr([]string{"apple", "bee", "cat"}, "bee")
	expected := true

	if result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	} else {
		t.Logf("expected: %v, got: %v", expected, result)
	}
}

func Test_ContainsStr_ReturnFalseWhenStringNotFoundInSlice(t *testing.T) {
	result := ContainsStr([]string{"apple", "bee", "cat"}, "dog")
	expected := false

	if result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	} else {
		t.Logf("expected: %v, got: %v", expected, result)
	}
}
