package main

import (
	"github.com/nugget/version-checkers/pkg/zonedirector"

	"encoding/json"
	"fmt"
)

func main() {
	releaseList, err := zonedirector.GetReleaseList(73)
	// fmt.Printf("There are %d releases in the list\n", len(releaseList))

	latest := zonedirector.FindLatest(releaseList, "IMG")
	// fmt.Printf("%+v\n", latest)

	j, err := json.Marshal(latest)
	if err != nil {
		panic(fmt.Sprintf("Cannot parse json: %v", err))
	}
	fmt.Println(string(j))
}
