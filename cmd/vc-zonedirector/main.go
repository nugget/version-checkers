package main

import (
	"github.com/nugget/version-checkers/internal/vc"
	"github.com/nugget/version-checkers/pkg/zonedirector"

	"encoding/json"
	"fmt"
)

func main() {
	var (
		productID int = 73 // Default 73 "zonedirector 1200"
	)

	vc.Init()
	state, err := vc.LoadState()
	if err != nil {
		panic(fmt.Sprintf("Cannot load state: %v", err))
	}

	fmt.Println("State loaded from %v", state.Timestamp)

	releaseList, err := zonedirector.GetReleaseList(productID)
	fmt.Printf("There are %d releases in the list\n", len(releaseList))

	latest := zonedirector.FindLatest(releaseList, "IMG")
	fmt.Printf("%+v\n", latest)

	json, err := json.Marshal(latest)
	if err != nil {
		panic(fmt.Sprintf("Cannot parse json: %v", err))
	}

	if state.IsNewer(latest.Version) {
		fmt.Printf("NEW VERSION %s (current %s)\n", latest.Version, state.CurrentVersion)

		state.RawVersion = json
		state.CurrentVersion = latest.Version

		err := state.Update()
		if err != nil {
			panic(fmt.Sprintf("Cannot update version: %v", err))
		}
	} else {
		fmt.Printf("boring version %s (current %s)\n", latest.Version, state.CurrentVersion)
	}

}
