package main

import (
	"encoding/xml"
	"fmt"
	"log"
)

type (
	product struct {
		XMLName xml.Name `xml:"Product"`
		Color   string   `xml:"Color"`
	}

	document struct {
		XMLName  xml.Name  `xml:"Products"`
		Products []product `xml:"Product"`
	}
)

const (
	docsNumber  = 5000
	lookUpColor = "White"
)

func main() {
	docs := make([]string, docsNumber)
	for i := range docs {
		docs[i] = fmt.Sprintf("data-%.4d.xml", i)
	}

	n := freq(lookUpColor, docs)

	log.Printf("Searching through %d files with products. Found products with color: %s %d times.", len(docs), lookUpColor, n)
}

func freq(color string, docs []string) int {
	var found int

	for _, doc := range docs {
		fileName := fmt.Sprintf("%s.xml", doc[:4])
		_ = fileName
	}

	return found
}
