package util

import (
	"golang.org/x/exp/constraints"
)

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func SliceToMap[K comparable, T any](slice []T, keySelector func(int, T) K) map[K]T {
	mapped := map[K]T{}
	for i, item := range slice {
		mapped[keySelector(i, item)] = item
	}
	return mapped
}
