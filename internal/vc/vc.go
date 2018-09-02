package vc

import (
	"fmt"
	"time"

	"github.com/mcuadros/go-version"
)

type State struct {
	Timestamp      time.Time
	CurrentVersion string
	RawVersion     []byte
}

func NewState() State {
	s := State{}

	return s
}

func (s *State) IsNewer(latestVersion string) bool {
	latestVersion = version.Normalize(latestVersion)

	fmt.Println("IsNewer:", s.CurrentVersion, latestVersion)

	return version.Compare(s.CurrentVersion, latestVersion, "<")
}

func (s *State) Update() error {
	s.Timestamp = time.Now()
	return nil
}

func Init() {
	fmt.Println("vc-data")
}

func LoadState() (State, error) {
	s := NewState()

	s.CurrentVersion = "10.0.0"

	return s, nil
}
