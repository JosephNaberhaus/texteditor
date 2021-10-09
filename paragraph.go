package go_text_editor

import (
	"bytes"
	"strings"
)

type paragraph []graphemeCluster

func (p paragraph) String() string {
	sb := strings.Builder{}
	for _, gc := range p {
		sb.Write(gc)
	}
	return sb.String()
}

func (p paragraph) Equal(p2 paragraph) bool {
	if len(p) != len(p2) {
		return false
	}

	for i := range p {
		if !bytes.Equal(p[i], p2[i]) {
			return false
		}
	}

	return true
}

func (t *TextEditor) Write(text string) {
	clusters := splitGraphemeClusters(text)
	for _, cluster := range clusters {
		t.writeCluster(cluster)
	}
}

func (t *TextEditor) writeCluster(cluster graphemeCluster) {
	if cluster.isNewline() {
		t.Newline()
	} else {
		// If the cursor is at the end of the paragraph then just append
		// Otherwise, we need to insert the new cluster in the middle of a paragraph
		if t.CursorIsAtEndOfParagraph() {
			t.paragraphs[t.cursorParagraph] = append(t.paragraphs[t.cursorParagraph], cluster)
		} else {
			left := t.CurParagraph()[:t.cursorPos]
			middle := paragraph{cluster}
			right := t.CurParagraph()[t.cursorPos:]
			t.paragraphs[t.cursorParagraph] = append(left, append(middle, right...)...)
		}

		t.Right()
	}
}

func (t *TextEditor) Backspace() {
	if t.CursorIsAtStartOfParagraph() {
		if t.CursorIsOnFirstParagraph() {
			return
		}

		curParagraph := t.CurParagraph()

		if t.CursorIsOnLastParagraph() {
			t.paragraphs = t.paragraphs[:t.cursorParagraph]
		} else {
			before := t.paragraphs[:t.cursorParagraph]
			after := t.paragraphs[t.cursorParagraph+1:]
			t.paragraphs = append(before, after...)
		}
		t.cursorParagraph--

		t.paragraphs[t.cursorParagraph] = append(t.paragraphs[t.cursorParagraph], curParagraph...)

		t.cursorPos = t.CurParagraphLength()
		t.cursorPreferredColumn = t.CursorColumn()
	} else {
		if t.CursorIsAtEndOfParagraph() {
			t.paragraphs[t.cursorParagraph] = t.CurParagraph()[:t.cursorPos-1]
		} else {
			left := t.CurParagraph()[:t.cursorPos-1]
			right := t.CurParagraph()[t.cursorPos:]
			t.paragraphs[t.cursorParagraph] = append(left, right...)
		}

		t.cursorPos--
	}
}

func (t *TextEditor) Newline() {
	// Text to the left of the cursor will remain. Text to the right of the paragraph will be moved to the new one
	left := t.CurParagraph()[:t.cursorPos]
	right := t.CurParagraph()[t.cursorPos:]

	if t.cursorParagraph == len(t.paragraphs) - 1 {
		t.paragraphs = append(t.paragraphs, make([]graphemeCluster, 0))
	} else {
		before := t.paragraphs[:t.cursorParagraph]
		after := t.paragraphs[t.cursorParagraph:]
		t.paragraphs = append(before, make([]graphemeCluster, 0))
		t.paragraphs = append(t.paragraphs, after...)
	}

	t.paragraphs[t.cursorParagraph] = left
	t.paragraphs[t.cursorParagraph + 1] = append(t.paragraphs[t.cursorParagraph + 1], right...)

	t.cursorParagraph++
	t.cursorPos = 0
	t.cursorPreferredColumn = 0
}
