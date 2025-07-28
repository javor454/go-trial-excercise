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

func BenchmarkFreq(b *testing.B) {
	docs := make([]string, docsNumber)
	for i := range docs {
		docs[i] = fmt.Sprintf("data-%.4d.xml", i)
	}

	b.ResetTimer() // Reset timer after setup
	for i := 0; i < b.N; i++ {
		freq(lookUpColor, docs)
	}
}
