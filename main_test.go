package main

import (
	"fmt"
	"testing"
)

const (
	whiteOccurencesSingleDoc = 4
)

func TestFreq(t *testing.T) {
	docs := make([]string, docsNumber)
	for i := range docs {
		docs[i] = fmt.Sprintf("data-%.4d.xml", i)
	}

	n := freq(lookUpColor, docs)

	expected := whiteOccurencesSingleDoc * docsNumber
	if n != expected {
		t.Errorf("expected %d occurences got %d", expected, n)
	}
}
