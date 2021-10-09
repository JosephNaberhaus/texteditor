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

func (p paragraph) Wrap(width int) []paragraph {
	if len(p) == 0 {
		return []paragraph{{}}
	}

	wrapped := make([]paragraph, 0, (len(p) / width) + 1)
	for i := 0; i < len(p); i += width {
		line := p[i:min(i + width, len(p))]
		wrapped = append(wrapped, line)
	}
	return wrapped
}

func (t *TextEditor) setParagraph(i int, para paragraph) {
	t.paragraphs[i] = para
	if i < len(t.wrappedLinesCache) {
		t.wrappedLinesCache[i] = nil
	}
}

func (t *TextEditor) setParagraphs(paras []paragraph) {
	t.paragraphs = paras
	t.wrappedLinesCache = nil
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
			t.setParagraph(t.cursorParagraph, append(t.paragraphs[t.cursorParagraph], cluster))
		} else {
			left := t.CurParagraph()[:t.cursorPos]
			middle := paragraph{cluster}
			right := t.CurParagraph()[t.cursorPos:]
			t.setParagraph(t.cursorParagraph, append(left, append(middle, right...)...))
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
			t.setParagraphs(t.paragraphs[:t.cursorParagraph])
		} else {
			before := t.paragraphs[:t.cursorParagraph]
			after := t.paragraphs[t.cursorParagraph+1:]
			t.setParagraphs(append(before, after...))
		}
		t.cursorParagraph--

		t.cursorPos = t.CurParagraphLength()
		t.cursorPreferredColumn = t.CursorColumn()

		t.setParagraph(t.cursorParagraph, append(t.paragraphs[t.cursorParagraph], curParagraph...))
	} else {
		if t.CursorIsAtEndOfParagraph() {
			t.setParagraph(t.cursorParagraph, t.CurParagraph()[:t.cursorPos-1])
		} else {
			left := t.CurParagraph()[:t.cursorPos-1]
			right := t.CurParagraph()[t.cursorPos:]
			t.setParagraph(t.cursorParagraph, append(left, right...))
		}

		t.cursorPos--
	}

	t.cursorPreferredColumn = t.CursorColumn()
}

func (t *TextEditor) Newline() {
	// Text to the left of the cursor will remain. Text to the right of the paragraph will be moved to the new one
	left := t.CurParagraph()[:t.cursorPos]
	right := t.CurParagraph()[t.cursorPos:]

	if t.cursorParagraph == len(t.paragraphs) - 1 {
		t.setParagraphs(append(t.paragraphs, make([]graphemeCluster, 0)))
	} else {
		before := t.paragraphs[:t.cursorParagraph]
		after := t.paragraphs[t.cursorParagraph:]
		t.setParagraphs(append(before, make([]graphemeCluster, 0)))
		t.setParagraphs(append(t.paragraphs, after...))
	}

	t.setParagraph(t.cursorParagraph, left)
	t.setParagraph(t.cursorParagraph + 1, append(t.paragraphs[t.cursorParagraph + 1], right...))

	t.cursorParagraph++
	t.cursorPos = 0
	t.cursorPreferredColumn = 0
}

func (t *TextEditor) wrapParagraphs() {
	if t.wrappedLinesCache == nil {
		t.wrappedLinesCache = make([][]paragraph, 0, len(t.paragraphs))
	}

	for i, para := range t.paragraphs {
		if i >= len(t.wrappedLinesCache) {
			t.wrappedLinesCache = append(t.wrappedLinesCache, nil)
		}

		if t.wrappedLinesCache[i] == nil {
			t.wrappedLinesCache[i] = para.Wrap(t.width)
		}
	}
}

func (t *TextEditor) String() string {
	t.wrapParagraphs()

	sb := strings.Builder{}
	for i, paras := range t.wrappedLinesCache {
		for j, para := range paras {
			sb.WriteString(para.String())

			if j != (len(paras) - 1) {
				sb.WriteRune('\n')
			}
		}

		if i != (len(t.wrappedLinesCache) - 1) {
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}