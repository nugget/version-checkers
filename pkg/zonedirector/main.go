package zonedirector

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mcuadros/go-version"
	"gopkg.in/headzoo/surf.v1"
)

var (
	BaseURI string = "https://support.ruckuswireless.com"
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

func (r *Release) ImportTD(s *goquery.Selection) error {
	r.BaseVersion, _ = s.Attr("data-version")
	r.ID, _ = s.Attr("id")

	s.Find("a").Each(func(i int, a *goquery.Selection) {
		r.URI, _ = a.Attr("href")
		r.URI = BaseURI + r.URI
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

	r.VersionString = ExpandedVersion(r.Version)

	return nil
}

func ExpandedVersion(version string) string {
	i := strings.Split(version, ".")
	if len(i) == 5 {
		return fmt.Sprintf("%s.%s.%s.%s build %s", i[0], i[1], i[2], i[3], i[4])
	}

	return version
}

func GetReleaseList(filter int) (releaseList []Release, err error) {
	bow := surf.NewBrowser()
	err = bow.Open(BaseURI + fmt.Sprintf("/software?filter=%d#firmwares", filter))
	if err != nil {
		return nil, err
	}

	bow.Find("tr.software").Each(func(i int, s *goquery.Selection) {
		r := Release{}
		r.ImportTD(s)

		releaseList = append(releaseList, r)
	})

	return releaseList, nil
}

func FindLatest(releaseList []Release, fileType string) Release {
	latest := Release{}

	for _, e := range releaseList {
		if fileType == "" || e.FileType == fileType {
			newer := version.Compare(e.Version, latest.Version, ">")

			if newer {
				latest = e
				// fmt.Printf("* %s (%s)\n", e.Description, e.Version)
			} else {
				// fmt.Printf("  %s (%s)\n", e.Description, e.Version)
			}
		}
	}

	return latest
}
