package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
)

var srcURL = "https://raw.githubusercontent.com/jpdillingham/slack-desktop-dark-theme/master/loader.js"
var space = func() { fmt.Println() }

func main() {
	// parse args
	srcURLInputPtr := flag.String("src", "", "The source of the patch code")
	flag.Parse()

	space()

	// figure out where Slack should be installed (OS specific) and make sure the directory exists
	path := getSlackDir()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Sprintf("Slack directory was expected to be %s, but was not found.", path))
	} else {
		fmt.Println(fmt.Sprintf("Slack found in %s.", path))
	}

	space()

	// list all of the installed Slack versions by matching subdirectories against 'app-N.N.N'
	fmt.Println("Installed Slack versions:")

	installedVersions := getInstalledSlackVersions(path)

	for i, f := range installedVersions {
		fmt.Printf("%d) %s\n", i, f)
	}

	fmt.Printf("\nPlease select the version to patch: ")

	if *srcURLInputPtr != "" {
		srcURL = *srcURLInputPtr
		fmt.Print(fmt.Sprintf("Fetching source from %s...", srcURL))
	} else {
		fmt.Print(fmt.Sprintf("Fetching source from default location %s...", srcURL))
	}

	getPatchSrc(srcURL)

	fmt.Print(" done.")
}

func getSlackDir() (path string) {
	home, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf("%s\\AppData\\Local\\slack", home)
	default:
		panic("unsupported OS.")
	}
}

func getInstalledSlackVersions(path string) (versions []string) {
	var installedVersions []string
	installedVersions = make([]string, 0)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic("error listing files in Slack directory.  Try again with escalated privileges?")
	}

	for _, f := range files {
		name := f.Name()
		isAppDir, _ := regexp.MatchString("^app-((\\d+\\.?){3})$", name)

		if f.IsDir() && isAppDir {
			installedVersions = append(installedVersions, f.Name())
		}
	}

	return installedVersions
}

func getPatchSrc(srcUrl string) (src string) {
	resp, err := http.Get(srcUrl)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}
