package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractVideoID(t *testing.T) {
	testCases := []struct {
		name     string
		fileName string
		expected string
	}{
		{
			name:     "Case 1",
			fileName: "1a8zuTfZhK0+=IS ＂STARVATION MODE＂ A REAL THING？ (What The Science Says).en.srt",
			expected: "1a8zuTfZhK0+",
		},
		{
			name:     "Case 2",
			fileName: "0a_fVS2s4Ho+=The Most Effective Way to Train HAMSTRINGS ｜ Training Science Explained.en.srt",
			expected: "0a_fVS2s4Ho+",
		},
		{
			name:     "empty filename",
			fileName: "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := extractVideoID(tc.fileName)
			assert.Equal(t, tc.expected, result, "extractVideoID(%q) = %q; want %q", tc.fileName, result, tc.expected)
		})
	}
}
