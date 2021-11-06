package go_text_editor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGraphemeCluster_isIn(t *testing.T) {
	assert.False(t, toGrapheme("a").isIn([]graphemeCluster{}))

	group := []graphemeCluster{
		toGrapheme("a"),
		toGrapheme("🚀"),
	}

	assert.True(t, toGrapheme("a").isIn(group))
	assert.True(t, toGrapheme("🚀").isIn(group))
	assert.False(t, toGrapheme("b").isIn(group))
}

func TestGraphemeCluster_isNewline(t *testing.T) {
	assert.True(t, toGrapheme("\n").isNewline())
	assert.True(t, toGrapheme("\r").isNewline())
	assert.True(t, toGrapheme("\r\n").isNewline())
	assert.False(t, toGrapheme(" ").isNewline())
}

func TestGraphemeCluster_String(t *testing.T) {
	assert.Equal(t, toGrapheme("a").String(), "a")
	assert.Equal(t, toGrapheme("🚀").String(), "🚀")
	assert.Equal(t, toGrapheme("👨‍👨‍👧‍👦").String(), "👨‍👨‍👧‍👦")
}

func TestSplitGraphemeClusters(t *testing.T) {
	expected := paragraph{
		toGrapheme("a"),
		toGrapheme("b"),
		toGrapheme("c"),
		toGrapheme(" "),
		toGrapheme("🚀"),
		toGrapheme("👨‍👨‍👧‍👦"),
		toGrapheme("\n"),
		toGrapheme("\r"),
		toGrapheme("\r\n"),
		toGrapheme("​"), // zero width space
	}

	assert.Equal(t, expected, splitGraphemeClusters("abc 🚀👨‍👨‍👧‍👦\n\r\r\n​"))
}

func TestToGrapheme(t *testing.T) {
	assert.Equal(t, graphemeCluster([]byte("a")), toGrapheme("a"))
	assert.Panics(t, func() {
		toGrapheme("more than one grapheme")
	})
}
