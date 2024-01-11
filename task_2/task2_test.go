package main

import (
	"testing"
)

func TestFormatString(t *testing.T) {
	testCases := []struct {
		s        string
		expected string
	}{
		{
			s:        "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			s:        "abcd",
			expected: "abcd",
		},
		{
			s:        "45",
			expected: "",
		},
		{
			s:        "",
			expected: "",
		},
		{
			s:        `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			s:        `qwe\45`,
			expected: `qwe44444`,
		},
		{
			s:        `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			s:        `\`,
			expected: "",
		},
		{
			s:        `qwe\\45`,
			expected: "",
		},
		{
			s:        `\\`,
			expected: `\`,
		},
		{
			s:        "1",
			expected: "",
		},
		{
			s:        `abc\`,
			expected: `abc`,
		},
		{
			s:        "a1b1b2",
			expected: "abbb",
		},
		{
			s:        `😀2😃😄😁`,
			expected: `😀😀😃😄😁`,
		},
	}
	for _, testCase := range testCases {
		res, err := FormatString(testCase.s)
		if res != testCase.expected {
			t.Error(err)
		}
	}
}
