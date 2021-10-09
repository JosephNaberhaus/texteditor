package go_text_editor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertParagraphsEqual(t *testing.T, expected, actual []paragraph) {
	if len(expected) != len(actual) {
		t.Errorf("Expected %d paragraphs, but got %d paragraphs", len(expected), len(actual))
	}

	for i := range expected {
		if i >= len(expected) || i >= len(actual) {
			return
		}

		if !expected[i].Equal(actual[i]) {
			message := fmt.Sprintf("Paragraphs at index %d are not equal:\nexpected: \"%s\"\nactual: \"%s\"", i, expected[i], actual[i])
			assert.Fail(t, message)
		}
	}
}
