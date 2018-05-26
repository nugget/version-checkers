package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mcuadros/go-version"
	"gopkg.in/headzoo/surf.v1"
)

type Release struct {
	BaseVersion   string
	ID            string
	Version       string
	VersionString string
	URI           string
	Description   string
	FileType      string
	Date          time.Time
}

func parseVersion(version string) string {
	i := strings.Split(version, ".")
	if len(i) == 5 {
		return fmt.Sprintf("%s.%s.%s.%s build %s", i[0], i[1], i[2], i[3], i[4])
	}

	return version
}

func main() {
	baseURI := "https://support.ruckuswireless.com"

	bow := surf.NewBrowser()
	err := bow.Open(baseURI + "/software?filter=73#firmwares")
	if err != nil {
		panic(err)
	}

	// fmt.Println(bow.Title())

	var releaseList []Release

	bow.Find("tr.software").Each(func(i int, s *goquery.Selection) {
		r := Release{}

		r.BaseVersion, _ = s.Attr("data-version")
		r.ID, _ = s.Attr("id")

		s.Find("a").Each(func(i int, a *goquery.Selection) {
			r.URI, _ = a.Attr("href")
			r.URI = baseURI + r.URI
		})

		s.Find("td").Each(func(i int, td *goquery.Selection) {
			switch i {
			case 0:
				r.Description = td.Text()
			case 1:
				r.Version = td.Text()
			case 2:
				r.FileType = td.Text()
			case 3:
				tt, err := time.Parse("2006-01-02", td.Text())
				if err != nil {
					panic(fmt.Sprintf("Error parsing date %v: %v\n", td.Text(), err))
				}
				r.Date = tt
			}
		})

		r.VersionString = parseVersion(r.Version)

		// fmt.Printf("%+v\n", r)
		releaseList = append(releaseList, r)
	})

	// fmt.Printf("There are %d releases in the list\n", len(releaseList))

	latest := Release{}

	for _, e := range releaseList {
		if e.FileType == "IMG" {
			newer := version.Compare(e.Version, latest.Version, ">")

			if newer {
				latest = e
				// fmt.Printf("* %s (%s)\n", e.Description, e.Version)
			} else {
				// fmt.Printf("  %s (%s)\n", e.Description, e.Version)
			}
		}
	}

	j, err := json.Marshal(latest)
	if err != nil {
		panic(fmt.Sprintf("Cannot parse json: %v", err))
	}
	fmt.Println(string(j))
}
