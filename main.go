package main

import (
	"bufio"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
)

// ---- Generic Interface ----
type CSVConvertible interface {
	ToCSV() []string
}

// ---- Structs ----

// Tags.xml
type Tag struct {
	ID            int    `xml:"Id,attr"`
	TagName       string `xml:"TagName,attr"`
	Count         int    `xml:"Count,attr"`
	ExcerptPostId *int   `xml:"ExcerptPostId,attr"`
	WikiPostId    *int   `xml:"WikiPostId,attr"`
}

func (t *Tag) ToCSV() []string {
	excerpt := ""
	if t.ExcerptPostId != nil {
		excerpt = strconv.Itoa(*t.ExcerptPostId)
	}
	wiki := ""
	if t.WikiPostId != nil {
		wiki = strconv.Itoa(*t.WikiPostId)
	}
	return []string{
		strconv.Itoa(t.ID),
		t.TagName,
		strconv.Itoa(t.Count),
		excerpt,
		wiki,
	}
}

// ---- Example Existing Structs (Post/User/Comment/Vote) ----
// (your previous code here… keeping them as-is)

// ---- Generic Converter ----
func convertXMLToCSV[T CSVConvertible](xmlPath, csvPath string, headers []string, newItem func() T) error {
	// input
	inFile, err := os.Open(xmlPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	decoder := xml.NewDecoder(bufio.NewReaderSize(inFile, 4<<20)) // 4 MB buffer

	// output
	outFile, err := os.Create(csvPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := csv.NewWriter(bufio.NewWriterSize(outFile, 4<<20)) // buffered writer
	defer writer.Flush()

	// header row
	if err := writer.Write(headers); err != nil {
		return err
	}

	// batching
	batchSize := 10000
	batch := make([][]string, 0, batchSize)

	// streaming parse
	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}
		switch se := tok.(type) {
		case xml.StartElement:
			if se.Name.Local == "row" {
				item := newItem() // allocate struct
				if err := decoder.DecodeElement(item, &se); err != nil {
					fmt.Println("decode error:", err)
					continue
				}
				batch = append(batch, item.ToCSV())
				if len(batch) == cap(batch) {
					if err := writer.WriteAll(batch); err != nil {
						return err
					}
					batch = batch[:0]
				}
			}
		}
	}

	// flush remaining
	if len(batch) > 0 {
		if err := writer.WriteAll(batch); err != nil {
			return err
		}
	}

	return nil
}

// ---- Example Usage ----
func main() {
	// Tags
	if err := convertXMLToCSV(
		"Tags.xml",
		"tags.csv",
		[]string{"id", "tag_name", "count", "excerpt_post_id", "wiki_post_id"},
		func() *Tag { return &Tag{} },
	); err != nil {
		panic(err)
	}
}
