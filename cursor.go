package editor

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

func (t *TextEditor) CursorIndex() int {
	return t.cursorPos - t.minCursorPos()
}

func (t *TextEditor) CursorColumn() int {
	return t.cursorPos % t.width
}

func (t *TextEditor) CursorRow() int {
	t.wrapParagraphs()

	row := t.cursorPos / t.width
	for i := 0; i < t.cursorParagraph; i++ {
		row += len(t.wrappedLinesCache[i])
	}

	return row
}

func (t *TextEditor) CursorIsAtStartOfParagraph() bool {
	return t.cursorPos == t.minCursorPos()
}

func (t *TextEditor) CursorIsAtEndOfParagraph() bool {
	return t.cursorPos == t.CurParagraphLength()
}

func (t *TextEditor) Left() {
	t.cursorPos--

	if t.cursorPos < t.minCursorPos() {
		if t.CursorIsOnFirstParagraph() {
			t.cursorPos = t.minCursorPos()
		} else {
			t.cursorParagraph--
			t.cursorPos = t.CurParagraphLength()
		}
	}

	t.cursorPreferredColumn = t.CursorColumn()
}

func (t *TextEditor) LeftNum(n int) {
	callNum(t.Left, n)
}

func (t *TextEditor) Right() {
	t.cursorPos++

	if t.cursorPos > (t.CurParagraphLength() + t.minCursorPos()) {
		if t.CursorIsOnLastParagraph() {
			t.cursorPos = t.CurParagraphLength() + t.minCursorPos()
		} else {
			t.cursorParagraph++
			t.cursorPos = t.minCursorPos()
		}
	}

	t.cursorPreferredColumn = t.CursorColumn()
}

func (t *TextEditor) RightNum(n int) {
	callNum(t.Right, n)
}

func (t *TextEditor) Up() {
	if t.cursorPos >= t.width {
		t.cursorPos -= (t.CursorColumn() - t.cursorPreferredColumn) + t.width
		t.cursorPos = max(t.minCursorPos(), t.cursorPos)
	} else {
		if !t.CursorIsOnFirstParagraph() {
			t.cursorParagraph--

			lineOffset := t.CurParagraphLength() / t.width
			t.cursorPos = max(
				t.minCursorPos(),
				min(
					t.CurParagraphLength(),
					(lineOffset*t.width)+t.cursorPreferredColumn,
				),
			)
		}
	}
}

func (t *TextEditor) UpNum(n int) {
	callNum(t.Up, n)
}

func (t *TextEditor) Down() {
	lineOffset := t.CurParagraphLength() / t.width
	if t.cursorPos < (lineOffset * t.width) {
		t.cursorPos = min(t.cursorPos+t.width, t.CurParagraphLength())
	} else {
		if !t.CursorIsOnLastParagraph() {
			t.cursorParagraph++
			t.cursorPos = min(t.cursorPreferredColumn, t.CurParagraphLength())
		}
	}
}

func (t *TextEditor) DownNum(n int) {
	callNum(t.Down, n)
}

func (t *TextEditor) Home() {
	t.cursorPos = t.minCursorPos()
	t.setPreferredColumn()
}

func (t *TextEditor) End() {
	t.cursorPos = t.CurParagraphLength() + t.minCursorPos()
	t.setPreferredColumn()
}

func (t *TextEditor) setPreferredColumn() {
	t.cursorPreferredColumn = t.cursorPos
}

func (t *TextEditor) minCursorPos() int {
	if t.cursorParagraph == 0 {
		return t.firstLineIndent
	}

	return 0
}
