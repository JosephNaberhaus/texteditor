package editor

import (
	"bytes"
	"errors"
	"github.com/rivo/uniseg"
)

// A single unicode grapheme cluster.
type graphemeCluster []byte

var newlineGraphemeClusters = []graphemeCluster{
	[]byte("\n"),
	[]byte("\r"),
	[]byte("\r\n"),
}

func (g graphemeCluster) isIn(array []graphemeCluster) bool {
	for _, toCheck := range array {
		if bytes.Equal(g, toCheck) {
			return true
		}
	}

	return false
}

func (g graphemeCluster) isNewline() bool {
	return g.isIn(newlineGraphemeClusters)
}

func (g graphemeCluster) String() string {
	return string(g)
}

// splitGraphemeClusters Splits a string into grapheme clusters.
func splitGraphemeClusters(s string) paragraph {
	clusters := make([]graphemeCluster, 0)

	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		clusters = append(clusters, gr.Bytes())
	}

	return clusters
}

// toGrapheme Converts the given string into a single grapheme cluster. Panics if the string doesn't contain exactly
// one grapheme cluster
func toGrapheme(s string) graphemeCluster {
	clusters := splitGraphemeClusters(s)
	if len(clusters) != 1 {
		panic(errors.New("string contains more than one grapheme cluster"))
	}

	return clusters[0]
}
