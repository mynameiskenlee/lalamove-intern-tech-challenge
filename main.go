package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run
	sort.Sort(semver.Versions(releases)) //sort the version from the smallest to largest
	//remove all versions lower than the minVersion first
	for i := len(releases) - 1; i >= 0; i-- { //going from the largest version to the smaller version to reduce the number of step needed to complete
		if releases[i].LessThan(*minVersion) { //check if the version is less than the minVersion
			releases = releases[i+1 : len(releases)] //slice out the version which is greater than the minVersion
			break                                    //since we only need to do the same operation once, there is no need to continue the loop
		}
	}
	sliceVer := releases[0].Slice() //slice the version to major, minor and patch array and store it
	version := releases[0]          //store the current version
	for i := 0; i < len(releases); i++ {
		temp := releases[i].Slice() //slice the version to major, minor and patch array and store it to a temp variable
		if temp[0] > sliceVer[0] || temp[1] > sliceVer[1]{  //check if the major or minor version changed
			versionSlice = append([]*semver.Version{version}, versionSlice...) //store the previous version to array
			sliceVer = temp                                                    //update the sliceVer variable
			version = releases[i]                                              //change to a higher version
		}
		if temp[2] > sliceVer[2] {
			sliceVer[2] = temp[2] //update the sliceVer variable
			version = releases[i] //change the current version to a higher patch
		}
	}
	versionSlice = append([]*semver.Version{version}, versionSlice...) //since the last one haven't been stored to the array, now store it
	return versionSlice
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	releases, _, err := client.Repositories.ListReleases(ctx, "kubernetes", "kubernetes", opt)
	if err != nil {
		panic(err) // is this really a good way?
	}
	minVersion := semver.New("1.11.0")
	allReleases := make([]*semver.Version, len(releases))
	for i, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
		}
		allReleases[i] = semver.New(versionString)
	}
	versionSlice := LatestVersions(allReleases, minVersion)

	fmt.Printf("latest versions of kubernetes/kubernetes: %s", versionSlice)
}
