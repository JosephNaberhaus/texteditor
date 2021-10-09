package go_text_editor

func (t *TextEditor) CurParagraph() paragraph {
	return t.paragraphs[t.cursorParagraph]
}

func (t *TextEditor) CurParagraphLength() int {
	return len(t.CurParagraph())
}

func (t *TextEditor) CursorIsOnFirstParagraph() bool {
	return t.cursorParagraph == 0
}

func (t *TextEditor) CursorIsOnLastParagraph() bool {
	return t.cursorParagraph == (len(t.paragraphs) - 1)
}

func (t *TextEditor) CursorColumn() int {
	return t.cursorPos % t.width
}

func (t *TextEditor) CursorIsAtStartOfParagraph() bool {
	return t.cursorPos == 0
}

func (t *TextEditor) CursorIsAtEndOfParagraph() bool {
	return t.cursorPos == t.CurParagraphLength()
}

func (t *TextEditor) Left() {
	t.cursorPos--

	if t.cursorPos < 0 {
		if t.CursorIsOnFirstParagraph() {
			t.cursorPos = 0
		} else {
			t.cursorParagraph--
			t.cursorPos = t.CurParagraphLength()
		}
	}

	t.cursorPreferredColumn = t.CursorColumn()
}

func (t *TextEditor) LeftNum(n int) {
	for i := 0; i < n; i++ {
		t.Left()
	}
}

func (t *TextEditor) Right() {
	t.cursorPos++

	if t.cursorPos > (t.CurParagraphLength() + 1) {
		if t.CursorIsOnLastParagraph() {
			t.cursorPos = t.CurParagraphLength()
		} else {
			t.cursorParagraph++
			t.cursorPos = 0
		}
	}

	t.cursorPreferredColumn = t.CursorColumn()
}

func (t *TextEditor) RightNum(n int) {
	for i := 0; i < n; i++ {
		t.Right()
	}
}