package editor

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParagraph_String(t *testing.T) {
	assert.Equal(t, splitGraphemeClusters("testing").String(), "testing")
	assert.Equal(t, splitGraphemeClusters("🥞").String(), "🥞")
	assert.Equal(t, splitGraphemeClusters("👨‍👨‍👦‍👦").String(), "👨‍👨‍👦‍👦")
}

func TestParagraph_Equal(t *testing.T) {
	p1 := splitGraphemeClusters("testing")
	p2 := splitGraphemeClusters("testing")
	p3 := splitGraphemeClusters("test")
	p4 := splitGraphemeClusters("helping")

	assert.True(t, p1.Equal(p1))
	assert.True(t, p1.Equal(p2))
	assert.False(t, p1.Equal(p3))
	assert.False(t, p1.Equal(p4))
}

func TestTextEditor_Write(t *testing.T) {
	twenty := "20-grapheme-clusters"

	editor := NewEditor()

	editor.Write(twenty)

	assertParagraphsEqual(t, []paragraph{splitGraphemeClusters(twenty)}, editor.paragraphs)
	assert.Equal(t, 20, editor.cursorPos)
	assert.Equal(t, 20, editor.cursorPreferredColumn)

	// Writing 3 more sets of twenty puts us right at the default editor width
	editor.Write(strings.Repeat(twenty, 3))

	assertParagraphsEqual(t, []paragraph{splitGraphemeClusters(strings.Repeat(twenty, 4))}, editor.paragraphs)
	assert.Equal(t, 80, editor.cursorPos)
	assert.Equal(t, 0, editor.cursorPreferredColumn)

	// Writing 1 more character should wrap the cursor but still write to the same paragraph
	editor.Write("a")

	assertParagraphsEqual(t, []paragraph{splitGraphemeClusters(strings.Repeat(twenty, 4) + "a")}, editor.paragraphs)
	assert.Equal(t, 81, editor.cursorPos)
	assert.Equal(t, 1, editor.cursorPreferredColumn)
}

func TestTextEditor_Write_MidParagraph(t *testing.T) {
	editor := NewEditor()

	editor.Write("Hello World!")
	editor.LeftNum(7)
	editor.Write("!")

	assertParagraphsEqual(t, []paragraph{splitGraphemeClusters("Hello! World!")}, editor.paragraphs)
	assert.Equal(t, editor.cursorPos, 6)
	assert.Equal(t, editor.cursorPreferredColumn, 6)

	editor.LeftNum(6)
	editor.Write("!")

	assertParagraphsEqual(t, []paragraph{splitGraphemeClusters("!Hello! World!")}, editor.paragraphs)
	assert.Equal(t, editor.cursorPos, 1)
	assert.Equal(t, editor.cursorPreferredColumn, 1)
}

func TestTextEditor_Write_Newline(t *testing.T) {
	editor := NewEditor()

	editor.Write("\n")

	assert.Equal(t, []paragraph{{}, {}}, editor.paragraphs)

	// Write some text, move the cursor to the middle of it, and then insert a newline
	editor.Write("afterthisfour")
	editor.LeftNum(4)
	editor.Write("\n")

	expected := []paragraph{{}, splitGraphemeClusters("afterthis"), splitGraphemeClusters("four")}
	assertParagraphsEqual(t, expected, editor.paragraphs)

	editor.LeftNum(5)
	editor.Write("\n")

	expected = []paragraph{
		{},
		splitGraphemeClusters("after"),
		splitGraphemeClusters("this"),
		splitGraphemeClusters("four"),
	}
	assertParagraphsEqual(t, expected, editor.paragraphs)
}

func TestTextEditor_Write_Regroups(t *testing.T) {
	editor := NewEditor()

	editor.Write("👨")
	assert.Equal(t, 1, len(editor.paragraphs[0]))
	assertParagraphsEqual(t, []paragraph{splitGraphemeClusters("👨")}, editor.paragraphs)

	editor.Write("\u200d👩")
	assert.Equal(t, 1, len(editor.paragraphs[0]))
	assertParagraphsEqual(t, []paragraph{splitGraphemeClusters("👨‍👩")}, editor.paragraphs)

	editor.Write("\u200d👦")
	assert.Equal(t, 1, len(editor.paragraphs[0]))
	assertParagraphsEqual(t, []paragraph{splitGraphemeClusters("👨‍👩‍👦")}, editor.paragraphs)
}

func TestTextEditor_Backspace(t *testing.T) {
	editor := NewEditor()

	editor.Write("four")
	editor.Backspace()

	assertParagraphsEqual(t, []paragraph{splitGraphemeClusters("fou")}, editor.paragraphs)
	assert.Equal(t, 3, editor.cursorPos)
	assert.Equal(t, 3, editor.cursorPreferredColumn)

	editor.LeftNum(3)
	editor.Backspace()

	assertParagraphsEqual(t, []paragraph{splitGraphemeClusters("fou")}, editor.paragraphs)
	assert.Equal(t, 0, editor.cursorPos)

	editor.RightNum(2)
	editor.Backspace()

	assertParagraphsEqual(t, []paragraph{splitGraphemeClusters("fu")}, editor.paragraphs)
	assert.Equal(t, 1, editor.cursorPos)
	assert.Equal(t, 1, editor.cursorPreferredColumn)
}

func TestTextEditor_Backspace_MultipleParagraphs(t *testing.T) {
	editor := NewEditor()

	editor.Write("Hello\n")
	editor.Write("World\n")
	editor.Backspace()

	expected := []paragraph{splitGraphemeClusters("Hello"), splitGraphemeClusters("World")}
	assertParagraphsEqual(t, expected, editor.paragraphs)
	assert.Equal(t, 5, editor.cursorPos)
	assert.Equal(t, 5, editor.cursorPreferredColumn)

	editor.LeftNum(5)
	editor.Backspace()

	assertParagraphsEqual(t, []paragraph{splitGraphemeClusters("HelloWorld")}, editor.paragraphs)
	assert.Equal(t, 5, editor.cursorPos)
	assert.Equal(t, 5, editor.cursorPreferredColumn)
}

func TestTextEditor_Backspace_MultipleParagraphs2(t *testing.T) {
	editor := NewEditor()

	editor.Write("Hello\n")
	editor.Write("World\n")
	editor.Write("End")
	editor.LeftNum(9)
	editor.Backspace()

	expected := []paragraph{splitGraphemeClusters("HelloWorld"), splitGraphemeClusters("End")}
	assertParagraphsEqual(t, expected, editor.paragraphs)
}

func TestTextEditor_Backspace_FirstLineIndent(t *testing.T) {
	editor := NewEditor()
	editor.SetFirstLineIndent(5)

	editor.Backspace()
	assert.Equal(t, 5, editor.cursorPos)
}

func TestTextEditor_String(t *testing.T) {
	editor := NewEditor()

	editor.Write("Hello\n")
	editor.Write("World!\n")
	editor.Write("\n")
	assert.Equal(t, "Hello\nWorld!\n\n", editor.String())

	editor.SetWidth(3)
	assert.Equal(t, "Hel\nlo\nWor\nld!\n\n", editor.String())
}
