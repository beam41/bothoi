package util

import (
	"testing"
)

func Test_Min_ReturnSmallerValueInt(t *testing.T) {
	result := Min(3, 1)
	expected := 1

	if result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	} else {
		t.Logf("expected: %v, got: %v", expected, result)
	}
}

func Test_Min_ReturnSmallerValueFloat(t *testing.T) {
	result := Min(3.0, 1.0)
	expected := 1.0

	if result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	} else {
		t.Logf("expected: %v, got: %v", expected, result)
	}
}

func Test_Min_ReturnSmallerValueStr(t *testing.T) {
	result := Min("apple", "bee")
	expected := "apple"

	if result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	} else {
		t.Logf("expected: %v, got: %v", expected, result)
	}
}

func Test_Ternary_ReturnFirstArgumentWhenTrue(t *testing.T) {
	result := Ternary(true, 1, 2)
	expected := 1

	if result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	} else {
		t.Logf("expected: %v, got: %v", expected, result)
	}
}

func Test_Ternary_ReturnSecondArgumentWhenFalse(t *testing.T) {
	result := Ternary(false, 1, 2)
	expected := 2

	if result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	} else {
		t.Logf("expected: %v, got: %v", expected, result)
	}
}

func Test_Contains_ReturnTrueWhenFoundInSlice(t *testing.T) {
	result := Contains([]string{"apple", "bee", "cat"}, "bee")
	expected := true

	if result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	} else {
		t.Logf("expected: %v, got: %v", expected, result)
	}
}

func Test_Contains_ReturnFalseWhenNotFoundInSlice(t *testing.T) {
	result := Contains([]string{"apple", "bee", "cat"}, "dog")
	expected := false

	if result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	} else {
		t.Logf("expected: %v, got: %v", expected, result)
	}
}
