package zonedirector

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mcuadros/go-version"
	"github.com/soniah/gosnmp"
	"gopkg.in/headzoo/surf.v1"
)

// BaseURI is the web location of the Ruckus support site.
const BaseURI = "https://support.ruckuswireless.com"

// Release represents a single firmware release from the vendor.
type Release struct {
	OK             bool
	BaseVersion    string
	ID             string
	Version        string
	VersionString  string
	URI            string
	Description    string
	FileType       string
	Date           time.Time
	Epoch          int64
	RunningVersion string
	VersionMatch   bool
}

// ImportRow parses an HTML row from the website into a Release object.
func (r *Release) ImportRow(s *goquery.Selection) error {
	var ok bool

	r.OK = false

	r.BaseVersion, ok = s.Attr("data-version")
	if !ok {
		return errors.New("ImportRow Unable to parse selection for data-version")
	}

	r.ID, ok = s.Attr("id")
	if !ok {
		return errors.New("ImportRow Unable to parse selection for id")
	}

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
			if err == nil {
				r.Date = tt
				r.Epoch = tt.Unix()
			}
		}
	})

	r.VersionString = ExpandedVersion(r.Version)

	r.OK = true

	return nil
}

func (r *Release) GetRunningVersion(hostname string, port uint16, community string) error {
	oids := []string{"1.3.6.1.4.1.25053.1.2.1.1.1.1.18.0"}

	params := &gosnmp.GoSNMP{
		Target:    hostname,
		Port:      port,
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
	}

	{
		err := params.Connect()
		if err != nil {
			return err
		}
	}

	defer params.Conn.Close()

	{
		result, err := params.Get(oids)
		if err != nil {
			return err
		}

		for _, variable := range result.Variables {
			r.RunningVersion = string(variable.Value.([]byte))
			return nil
		}
	}

	return errors.New("Unable to find current version from appliance")
}

// ExpandedVersion takes a literal version string and converts it into tha
// human-readable form exposed in SNMP and in the ZoneDirector web console.
func ExpandedVersion(version string) string {
	i := strings.Split(version, ".")
	if len(i) == 5 {
		return fmt.Sprintf("%s.%s.%s.%s build %s", i[0], i[1], i[2], i[3], i[4])
	}

	return version
}

// GetReleaseList fetches the list of current firmware releases from the vendor website.
func GetReleaseList(filter int) (releaseList []Release, err error) {
	bow := surf.NewBrowser()
	err = bow.Open(BaseURI + fmt.Sprintf("/software?filter=%d#firmwares", filter))
	if err != nil {
		return nil, err
	}

	bow.Find("tr.software").Each(func(i int, s *goquery.Selection) {
		r := Release{}
		r.ImportRow(s)

		releaseList = append(releaseList, r)
	})

	return releaseList, nil
}

// FindLatest examines a release list slice and finds the latest version within.
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
