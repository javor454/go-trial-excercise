package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func freqSequentialV1(color string, docs []string) int {
	var found int

	for _, doc := range docs {
		fileName := fmt.Sprintf("%s.xml", doc[:4])

		bytes, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatal(err)
		}

		var doc document
		err = xml.Unmarshal(bytes, &doc)
		if err != nil {
			log.Fatal(err)
		}

		for _, product := range doc.Products {
			if strings.Contains(product.Color, color) {
				found++
			}
		}
	}

	return found
}

func freqSequentialV2(color string, docs []string) int {
	var found int

	for _, doc := range docs {
		fileName := fmt.Sprintf("%s.xml", doc[:4])

		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		decoder := xml.NewDecoder(file)
		for {
			token, err := decoder.Token()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				log.Fatal(err)
			}

			switch t := token.(type) {
			case xml.StartElement:
				if t.Name.Local == "Product" {
					var product product
					err = decoder.DecodeElement(&product, &t)
					if err != nil {
						log.Fatal(err)
					}

					if strings.Contains(product.Color, color) {
						found++
					}
				}
			}
		}
	}

	return found
}
