package solution

import "encoding/xml"

type (
	Product struct {
		XMLName xml.Name `xml:"Product"`
		Color   string   `xml:"Color"`
	}

	Document struct {
		XMLName  xml.Name  `xml:"Products"`
		Products []Product `xml:"Product"`
	}
)
