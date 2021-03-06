package hellfire // import "pathspider.net/hellfire"

import (
	"archive/zip"
	"log"
)

type AlexaTopsitesList struct {
	TestList
	filename string
}

// URL to download the latest Alexa Topsites list from
const AlexaTopsitesURL string = "http://s3.amazonaws.com/alexa-static/top-1m.csv.zip"

func (l *AlexaTopsitesList) SetFilename(filename string) {
	l.filename = filename
}

func (l *AlexaTopsitesList) FeedJobs(jobs chan map[string]interface{}) {
	var topsites *CSVList

	if l.filename == "" {
		urlReader, err := getReaderFromUrl(AlexaTopsitesURL)
		if err != nil {
			log.Fatalf("Unable to get <%s>: %s", AlexaTopsitesURL, err)
		}

		zr, err := zip.NewReader(urlReader, int64(urlReader.Len()))
		if err != nil {
			log.Fatalf("Unable to read zip: %s", err)
		}

		for _, zf := range zr.File {
			if zf.Name == "top-1m.csv" {
				f, _ := zf.Open()
				topsites = CSVListFromReader(f)
				break
			}
		}

		if topsites == nil {
			panic("Did not find top-1m.csv in the zip archive")
		}
	} else {
		topsites = CSVListFromFile(l.filename)
	}

	topsites.SetHeader([]string{"rank", "domain"})
	topsites.FeedJobs(jobs)
}
