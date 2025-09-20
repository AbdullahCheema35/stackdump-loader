package main

import (
	"bufio"
	"encoding/csv"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

// Users.xml
type User struct {
	ID              int    `xml:"Id,attr"`
	Reputation      int    `xml:"Reputation,attr"`
	CreationDate    string `xml:"CreationDate,attr"`
	DisplayName     string `xml:"DisplayName,attr"`
	LastAccessDate  string `xml:"LastAccessDate,attr"`
	WebsiteURL      string `xml:"WebsiteUrl,attr"`
	Location        string `xml:"Location,attr"`
	AboutMe         string `xml:"AboutMe,attr"`
	Views           int    `xml:"Views,attr"`
	UpVotes         int    `xml:"UpVotes,attr"`
	DownVotes       int    `xml:"DownVotes,attr"`
	ProfileImageURL string `xml:"ProfileImageUrl,attr"`
	AccountID       *int   `xml:"AccountId,attr"`
}

func (u *User) ToCSV() []string {
	account := ""
	if u.AccountID != nil {
		account = strconv.Itoa(*u.AccountID)
	}
	return []string{
		strconv.Itoa(u.ID),
		strconv.Itoa(u.Reputation),
		u.CreationDate,
		u.DisplayName,
		u.LastAccessDate,
		u.WebsiteURL,
		u.Location,
		u.AboutMe,
		strconv.Itoa(u.Views),
		strconv.Itoa(u.UpVotes),
		strconv.Itoa(u.DownVotes),
		u.ProfileImageURL,
		account,
	}
}

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
				item := newItem()
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

// ---- Main ----
func main() {
	inPath := flag.String("in", "", "Path to XML file (Tags.xml or Users.xml)")
	flag.Parse()

	if *inPath == "" {
		fmt.Println("Please provide -in=/path/to/Tags.xml or Users.xml")
		return
	}

	base := strings.ToLower(filepath.Base(*inPath))
	outPath := filepath.Join(filepath.Dir(*inPath), strings.TrimSuffix(filepath.Base(*inPath), filepath.Ext(*inPath))+".csv")

	switch base {
	case "tags.xml":
		if err := convertXMLToCSV(
			*inPath,
			outPath,
			[]string{"id", "tag_name", "count", "excerpt_post_id", "wiki_post_id"},
			func() *Tag { return &Tag{} },
		); err != nil {
			panic(err)
		}
	case "users.xml":
		if err := convertXMLToCSV(
			*inPath,
			outPath,
			[]string{"id", "reputation", "creation_date", "display_name", "last_access_date", "website_url", "location", "about_me", "views", "up_votes", "down_votes", "profile_image_url", "account_id"},
			func() *User { return &User{} },
		); err != nil {
			panic(err)
		}
	default:
		fmt.Printf("Unsupported file: %s (only Tags.xml and Users.xml supported)\n", base)
	}
}
