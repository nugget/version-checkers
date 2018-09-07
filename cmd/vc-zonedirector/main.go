package main

import (
	"github.com/nugget/version-checkers/pkg/zonedirector"

	"encoding/json"
	"fmt"
)

func main() {
	var (
		productID int = 73 // Default 73 "zonedirector 1200"
	)

	releaseList, err := zonedirector.GetReleaseList(productID)
	//fmt.Printf("There are %d releases in the list\n", len(releaseList))

	latest := zonedirector.FindLatest(releaseList, "IMG")
	//fmt.Printf("latest: %+v\n\n", latest)

	json, err := json.Marshal(latest)
	if err != nil {
		panic(fmt.Sprintf("Cannot parse json: %v", err))
	}

	fmt.Printf("%s\n", json)
}
