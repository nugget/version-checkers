package anthemav

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mcuadros/go-version"
	"gopkg.in/headzoo/surf.v1"
)

// BaseURI is the web location of the Ruckus support site.
const BaseURI = "https://www.anthemav.com/support/latest-software.php"

// Release represents a single software release from the vendor.
type Release struct {
	OK         bool
	ARCVersion string
	ARCURI     string
	FWVersion  string
	FWURI      string
}
