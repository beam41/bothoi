package util

import "testing"

func Test_ConvertVidLengthToSeconds_OnlySeconds(t *testing.T) {
	type matrix struct {
		test     string
		expected uint32
	}
	matrices := []matrix{
		{"0", 0},
		{"02", 2},
		{"10", 10},
		{"0:30", 30},
		{"0:00:20", 20},
	}

	for i, test := range matrices {
		result := ConvertVidLengthToSeconds(test.test)

		if result != test.expected {
			t.Errorf("%v: expected: %v, got: %v", i, test.expected, result)
		} else {
			t.Logf("%v: expected: %v, got: %v", i, test.expected, result)
		}
	}
}

func Test_ConvertVidLengthToSeconds_WithMinute(t *testing.T) {
	type matrix struct {
		test     string
		expected uint32
	}
	matrices := []matrix{
		{"00:00", 0},
		{"1:05", 65},
		{"05:30", 330},
		{"10:30", 630},
		{"70:30", 4230},
		{"0:1:20", 80},
		{"0:20:20", 1220},
	}

	for i, test := range matrices {
		result := ConvertVidLengthToSeconds(test.test)

		if result != test.expected {
			t.Errorf("%v: expected: %v, got: %v", i, test.expected, result)
		} else {
			t.Logf("%v: expected: %v, got: %v", i, test.expected, result)
		}
	}
}

func Test_ConvertVidLengthToSeconds_WithHours(t *testing.T) {
	type matrix struct {
		test     string
		expected uint32
	}
	matrices := []matrix{
		{"1:1:05", 3665},
		{"02:1:05", 7265},
		{"20:05:30", 72330},
		{"300:00:00", 1_080_000},
	}

	for i, test := range matrices {
		result := ConvertVidLengthToSeconds(test.test)

		if result != test.expected {
			t.Errorf("%v: expected: %v, got: %v", i, test.expected, result)
		} else {
			t.Logf("%v: expected: %v, got: %v", i, test.expected, result)
		}
	}
}

func Test_ConvertSecondsToVidLength_LessThanMinuteShouldHaveLeadingZero(t *testing.T) {
	type matrix struct {
		test     uint32
		expected string
	}
	matrices := []matrix{
		{5, "0:05"},
		{59, "0:59"},
	}

	for i, test := range matrices {
		result := ConvertSecondsToVidLength(test.test)

		if result != test.expected {
			t.Errorf("%v: expected: %v, got: %v", i, test.expected, result)
		} else {
			t.Logf("%v: expected: %v, got: %v", i, test.expected, result)
		}
	}
}

func Test_ConvertSecondsToVidLength_HaveTwoDigitWithoutLeadingZero(t *testing.T) {
	type matrix struct {
		test     uint32
		expected string
	}
	matrices := []matrix{
		{62, "1:02"},
		{742, "12:22"},
		{3742, "1:02:22"},
	}

	for i, test := range matrices {
		result := ConvertSecondsToVidLength(test.test)

		if result != test.expected {
			t.Errorf("%v: expected: %v, got: %v", i, test.expected, result)
		} else {
			t.Logf("%v: expected: %v, got: %v", i, test.expected, result)
		}
	}
}

func Test_ConvertSecondsToVidLength_HourCanHaveAsManyDigit(t *testing.T) {
	type matrix struct {
		test     uint32
		expected string
	}
	matrices := []matrix{
		{36000, "10:00:00"},
		{1152000, "320:00:00"},
		{11667600, "3241:00:00"},
	}

	for i, test := range matrices {
		result := ConvertSecondsToVidLength(test.test)

		if result != test.expected {
			t.Errorf("%v: expected: %v, got: %v", i, test.expected, result)
		} else {
			t.Logf("%v: expected: %v, got: %v", i, test.expected, result)
		}
	}
}
