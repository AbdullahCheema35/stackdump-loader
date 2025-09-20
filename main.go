package main

import (
	"bufio"
	"encoding/csv"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// ---- Generic Interface ----
type CSVConvertible interface {
	ToCSV() []string
}

// ---- Job Queue ----
type Job struct {
	FileName string
}

var jobs = make(chan Job, 100) // buffered channel
var done = make(chan struct{})

func worker(jobs <-chan Job, done chan<- struct{}, base string) {
	base = strings.ToLower(base) // e.g., tags, users, badges, votes, comments

	// process jobs
	for job := range jobs {
		fmt.Println("⚙️ Processing:", job.FileName)

		cmd := exec.Command("make", base, fmt.Sprintf("filename=%s", job.FileName))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ command failed for %s: %v\n", job.FileName, err)
			continue
		}

		// success → delete file
		if err := os.Remove(job.FileName); err != nil {
			fmt.Printf("⚠️ failed to delete %s: %v\n", job.FileName, err)
		} else {
			fmt.Printf("✅ success, deleted %s\n", job.FileName)
		}
	}

	close(done)
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

// Badges.xml
type Badge struct {
	ID       int    `xml:"Id,attr"`
	UserId   int    `xml:"UserId,attr"`
	Name     string `xml:"Name,attr"`
	Date     string `xml:"Date,attr"`
	Class    int    `xml:"Class,attr"`
	TagBased bool   `xml:"TagBased,attr"`
}

func (b *Badge) ToCSV() []string {
	return []string{
		strconv.Itoa(b.ID),
		strconv.Itoa(b.UserId),
		b.Name,
		b.Date,
		strconv.Itoa(b.Class),
		strconv.FormatBool(b.TagBased),
	}
}

// Votes.xml
type Vote struct {
	ID           int    `xml:"Id,attr"`
	PostId       int    `xml:"PostId,attr"`
	VoteTypeId   int    `xml:"VoteTypeId,attr"`
	UserId       *int   `xml:"UserId,attr"`
	CreationDate string `xml:"CreationDate,attr"`
	BountyAmount *int   `xml:"BountyAmount,attr"`
}

func (v *Vote) ToCSV() []string {
	user := ""
	if v.UserId != nil {
		user = strconv.Itoa(*v.UserId)
	}
	bounty := ""
	if v.BountyAmount != nil {
		bounty = strconv.Itoa(*v.BountyAmount)
	}
	return []string{
		strconv.Itoa(v.ID),
		strconv.Itoa(v.PostId),
		strconv.Itoa(v.VoteTypeId),
		user,
		v.CreationDate,
		bounty,
	}
}

// Comments.xml
type Comment struct {
	ID              int    `xml:"Id,attr"`
	PostId          int    `xml:"PostId,attr"`
	Score           int    `xml:"Score,attr"`
	Text            string `xml:"Text,attr"`
	CreationDate    string `xml:"CreationDate,attr"`
	UserDisplayName string `xml:"UserDisplayName,attr"`
	UserId          *int   `xml:"UserId,attr"`
	ContentLicense  string `xml:"ContentLicense,attr"`
}

func (c *Comment) ToCSV() []string {
	user := ""
	if c.UserId != nil {
		user = strconv.Itoa(*c.UserId)
	}

	// Escape COPY end marker `\.` so psql won't misinterpret it
	safeText := strings.ReplaceAll(c.Text, `\.`, `\\.`)

	return []string{
		strconv.Itoa(c.ID),
		strconv.Itoa(c.PostId),
		strconv.Itoa(c.Score),
		safeText,
		c.CreationDate,
		c.UserDisplayName,
		user,
		c.ContentLicense,
	}
}

// ---- Generic Converter with chunking ----
func convertXMLToCSV[T CSVConvertible](xmlPath, baseName string, headers []string, newItem func() T) error {
	// input
	inFile, err := os.Open(xmlPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	decoder := xml.NewDecoder(bufio.NewReaderSize(inFile, 4<<20)) // 4 MB buffer

	// chunking settings
	const chunkSize = 10_000_000 // 1 crore rows per chunk
	const batchSize = 10000      // write in batches of 10k rows
	fileIndex := 1
	recordCount := 0

	var outFile *os.File
	var writer *csv.Writer
	var bufferedWriter *bufio.Writer

	openNewFile := func() error {
		if outFile != nil {
			// finalize previous file
			writer.Flush()
			bufferedWriter.Flush()
			outFile.Close()

			// enqueue job
			fileName := fmt.Sprintf("%d_%s.csv", fileIndex, baseName)
			jobs <- Job{FileName: fileName}

			fileIndex++
		}

		// new file
		fileName := fmt.Sprintf("%d_%s.csv", fileIndex, baseName)
		var err error
		outFile, err = os.Create(fileName)
		if err != nil {
			return err
		}
		bufferedWriter = bufio.NewWriterSize(outFile, 4<<20)
		writer = csv.NewWriter(bufferedWriter)

		// write headers
		if err := writer.Write(headers); err != nil {
			return err
		}

		recordCount = 0
		return nil
	}

	if err := openNewFile(); err != nil {
		return err
	}

	// batching
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
				recordCount++

				if len(batch) == cap(batch) {
					if err := writer.WriteAll(batch); err != nil {
						return err
					}
					batch = batch[:0]
				}

				if recordCount >= chunkSize {
					// flush remaining in current batch
					if len(batch) > 0 {
						if err := writer.WriteAll(batch); err != nil {
							return err
						}
						batch = batch[:0]
					}
					// open next file
					if err := openNewFile(); err != nil {
						return err
					}
				}
			}
		}
	}

	// flush last batch
	if len(batch) > 0 {
		if err := writer.WriteAll(batch); err != nil {
			return err
		}
	}

	// finalize last file
	writer.Flush()
	bufferedWriter.Flush()
	outFile.Close()

	// enqueue job for last file
	fileName := fmt.Sprintf("%d_%s.csv", fileIndex, baseName)
	jobs <- Job{FileName: fileName}

	return nil
}

// ---- Main ----
func main() {
	inPath := flag.String("in", "", "Path to XML file (Tags.xml, Users.xml, Badges.xml, Votes.xml, Comments.xml)")
	flag.Parse()

	if *inPath == "" {
		fmt.Println("Please provide -in=/path/to/Tags.xml, Users.xml, Badges.xml, Votes.xml, or Comments.xml")
		return
	}

	base := strings.TrimSuffix(filepath.Base(*inPath), filepath.Ext(*inPath))

	// start worker
	go worker(jobs, done, base)

	switch strings.ToLower(filepath.Base(*inPath)) {
	case "tags.xml":
		if err := convertXMLToCSV(
			*inPath,
			base,
			[]string{"id", "tag_name", "count", "excerpt_post_id", "wiki_post_id"},
			func() *Tag { return &Tag{} },
		); err != nil {
			panic(err)
		}
	case "users.xml":
		if err := convertXMLToCSV(
			*inPath,
			base,
			[]string{"id", "reputation", "creation_date", "display_name", "last_access_date", "website_url", "location", "about_me", "views", "up_votes", "down_votes", "profile_image_url", "account_id"},
			func() *User { return &User{} },
		); err != nil {
			panic(err)
		}
	case "badges.xml":
		if err := convertXMLToCSV(
			*inPath,
			base,
			[]string{"id", "user_id", "name", "date", "class", "tag_based"},
			func() *Badge { return &Badge{} },
		); err != nil {
			panic(err)
		}
	case "votes.xml":
		if err := convertXMLToCSV(
			*inPath,
			base,
			[]string{"id", "post_id", "vote_type_id", "user_id", "creation_date", "bounty_amount"},
			func() *Vote { return &Vote{} },
		); err != nil {
			panic(err)
		}
	case "comments.xml":
		if err := convertXMLToCSV(
			*inPath,
			base,
			[]string{"id", "post_id", "score", "text", "creation_date", "user_display_name", "user_id", "content_license"},
			func() *Comment { return &Comment{} },
		); err != nil {
			panic(err)
		}
	default:
		fmt.Printf("Unsupported file: %s\n", *inPath)
	}

	// no more jobs → close channel
	close(jobs)

	// wait for worker
	<-done
	fmt.Println("✅ All jobs processed, exiting.")
}
