package go_text_editor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTextEditor_CurParagraph(t *testing.T) {
	editor := NewEditor()

	assert.Equal(t, editor.paragraphs[0], editor.CurParagraph())

	editor.Newline()
	editor.Write("testing")
	assert.Equal(t, editor.paragraphs[1], editor.CurParagraph())

	editor.Up()
	assert.Equal(t, editor.paragraphs[0], editor.CurParagraph())
}

func TestTextEditor_CurParagraphLength(t *testing.T) {
	editor := NewEditor()

	assert.Zero(t, editor.CurParagraphLength())

	editor.Write("four")
	assert.Equal(t, 4, editor.CurParagraphLength())

	editor.Newline()
	assert.Zero(t, editor.CurParagraphLength())
}

func TestTextEditor_CursorIsOnFirstParagraph(t *testing.T) {
	editor := NewEditor()

	assert.True(t, editor.CursorIsOnFirstParagraph())

	editor.Newline()
	assert.False(t, editor.CursorIsOnFirstParagraph())

	editor.Up()
	assert.True(t, editor.CursorIsOnFirstParagraph())
}

func TestTextEditor_CursorIsOnLastParagraph(t *testing.T) {
	editor := NewEditor()

	assert.True(t, editor.CursorIsOnLastParagraph())

	editor.Newline()
	assert.True(t, editor.CursorIsOnLastParagraph())

	editor.Up()
	assert.False(t, editor.CursorIsOnLastParagraph())
}

func TestTextEditor_CursorColumn(t *testing.T) {
	editor := NewEditor()
	editor.SetWidth(10)

	assert.Zero(t, editor.CursorColumn())

	editor.Write("fourfour")
	assert.Equal(t, 8, editor.CursorColumn())

	editor.Write("four")
	assert.Equal(t, 2, editor.CursorColumn())

	editor.Newline()
	assert.Zero(t, editor.CursorColumn())
}

func TestTextEditor_CursorRow(t *testing.T) {
	editor := NewEditor()
	editor.SetWidth(10)

	assert.Zero(t, editor.CursorRow())

	editor.Write("fourfour")
	assert.Zero(t, editor.CursorRow())

	editor.Write("four")
	assert.Equal(t, 1, editor.CursorRow())

	editor.SetWidth(100)
	assert.Zero(t, editor.CursorRow())

	editor.Newline()
	assert.Equal(t, 1, editor.CursorRow())
}

func TestTextEditor_CursorIsAtStartOfParagraph(t *testing.T) {
	editor := NewEditor()

	assert.True(t, editor.CursorIsAtStartOfParagraph())

	editor.Write("four")
	assert.False(t, editor.CursorIsAtStartOfParagraph())

	editor.LeftNum(4)
	assert.True(t, editor.CursorIsAtStartOfParagraph())

	editor.RightNum(4)
	editor.Newline()
	assert.True(t, editor.CursorIsAtStartOfParagraph())
}

func TestTextEditor_CursorIsAtEndOfParagraph(t *testing.T) {
	editor := NewEditor()

	assert.True(t, editor.CursorIsAtEndOfParagraph())

	editor.Write("four")
	assert.True(t, editor.CursorIsAtEndOfParagraph())

	editor.LeftNum(4)
	assert.False(t, editor.CursorIsAtEndOfParagraph())

	editor.RightNum(4)
	editor.Newline()
	assert.True(t, editor.CursorIsAtEndOfParagraph())
}

func TestTextEditor_Left(t *testing.T) {
	editor := NewEditor()

	editor.Left()
	assert.Zero(t, editor.cursorParagraph)
	assert.Zero(t, editor.cursorPos)

	editor.Write("four")
	editor.Left()
	assert.Zero(t, editor.cursorParagraph)
	assert.Equal(t, 3, editor.cursorPos)

	editor.Newline()
	editor.Write("ou")
	editor.Left()
	assert.Equal(t, 1, editor.cursorParagraph)
	assert.Equal(t, 1, editor.cursorPos)

	editor.Left()
	editor.Left()
	assert.Zero(t, editor.cursorParagraph)
	assert.Equal(t, 3, editor.cursorPos)
}

func TestTextEditor_LeftNum(t *testing.T) {
	editor := NewEditor()

	editor.LeftNum(10)
	assert.Zero(t, editor.cursorPos)

	editor.Write("fourfour")
	editor.LeftNum(4)
	assert.Equal(t, 4, editor.cursorPos)
}

func TestTextEditor_Right(t *testing.T) {
	editor := NewEditor()
	editor.Write("four")
	editor.Newline()
	editor.Write("or")
	editor.cursorParagraph = 0
	editor.cursorPos = 3

	editor.Right()
	assert.Zero(t, editor.cursorParagraph)
	assert.Equal(t, 4, editor.cursorPos)

	editor.Right()
	assert.Equal(t, 1, editor.cursorParagraph)
	assert.Equal(t, 0, editor.cursorPos)
}

func TestTextEditor_RightNum(t *testing.T) {
	editor := NewEditor()
	editor.Write("fourfour")
	editor.cursorPos = 0

	editor.RightNum(4)
	assert.Equal(t, 4, editor.cursorPos)

	editor.RightNum(4)
	assert.Equal(t, 8, editor.cursorPos)
}

func TestTextEditor_Up(t *testing.T) {
	editor := NewEditor()

	editor.Up()
	assert.Zero(t, editor.cursorParagraph)
	assert.Zero(t, editor.cursorPos)

	editor.Write("fourfourfour")
	editor.Newline()
	editor.Write("four")
	editor.Newline()
	editor.Write("fourfour")

	editor.Up()
	assert.Equal(t, 1, editor.cursorParagraph)
	assert.Equal(t, 4, editor.cursorPos)

	editor.Up()
	assert.Zero(t, editor.cursorParagraph)
	assert.Equal(t, 8, editor.cursorPos)
}

func TestTextEditor_Up_LongParagraph(t *testing.T) {
	editor := NewEditor()
	editor.SetWidth(8)

	editor.Write("fourfourfourfourfour")

	editor.Up()
	assert.Zero(t, editor.cursorParagraph)
	assert.Equal(t, 12, editor.cursorPos)
}

func TestTextEditor_UpNum(t *testing.T) {
	editor := NewEditor()

	editor.UpNum(8)
	assert.Zero(t, editor.cursorParagraph)
	assert.Zero(t, editor.cursorPos)

	editor.Newline()
	editor.Newline()
	editor.Newline()
	editor.Newline()
	editor.Newline()

	editor.UpNum(2)
	assert.Zero(t, editor.cursorPos)
	assert.Equal(t, 3, editor.cursorParagraph)

	editor.UpNum(3)
	assert.Zero(t, editor.cursorParagraph)
	assert.Zero(t, editor.cursorPos)
}

func TestTextEditor_Down(t *testing.T) {
	editor := NewEditor()

	editor.Down()
	assert.Zero(t, editor.cursorParagraph)
	assert.Zero(t, editor.cursorPos)

	editor.Write("fourfour")
	editor.Newline()
	editor.Write("four")
	editor.Newline()
	editor.Write("fourfourfour")
	editor.cursorPos = 8
	editor.cursorPreferredColumn = 8
	editor.cursorParagraph = 0

	editor.Down()
	assert.Equal(t, 1, editor.cursorParagraph)
	assert.Equal(t, 4, editor.cursorPos)

	editor.Down()
	assert.Equal(t, 2, editor.cursorParagraph)
	assert.Equal(t, 8, editor.cursorPos)
}

func TestTextEditor_Down_LongParagraph(t *testing.T) {
	editor := NewEditor()
	editor.SetWidth(8)

	editor.Write("fourfourfourfour")
	editor.cursorPos = 0
	editor.cursorParagraph = 0

	editor.Down()
	assert.Zero(t, editor.cursorParagraph)
	assert.Equal(t, 8, editor.cursorPos)
}

func TestTextEditor_DownNum(t *testing.T) {
	editor := NewEditor()

	editor.DownNum(8)
	assert.Zero(t, editor.cursorParagraph)
	assert.Zero(t, editor.cursorPos)

	editor.Newline()
	editor.Newline()
	editor.Newline()
	editor.Newline()
	editor.Newline()

	editor.cursorParagraph = 0
	editor.DownNum(2)
	assert.Equal(t, 2, editor.cursorParagraph)

	editor.DownNum(3)
	assert.Equal(t, 5, editor.cursorParagraph)
}

func TestTextEditor_Home(t *testing.T) {
	editor := NewEditor()
	editor.SetWidth(4)

	editor.Write("fourfourfourfour")
	editor.Home()
	assert.Zero(t, editor.cursorPos)
}

func TestTextEditor_End(t *testing.T) {
	editor := NewEditor()
	editor.SetWidth(4)

	editor.Write("fourfourfourfour")
	editor.cursorPos = 0
	editor.End()
	assert.Equal(t, 16, editor.cursorPos)
}
