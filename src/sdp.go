package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"

	util "./util"
)

var srcURL = "https://raw.githubusercontent.com/jpdillingham/slack-desktop-dark-theme/master/loader.js"
var space = func() { fmt.Println() }
var sep = os.PathSeparator

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

	space()

	// prompt user to select the target installation and assume they will enter something valid
	selectedVersionInput := util.PromptForInput("Please select the version to patch: ")
	selectedVersionIndex, _ := strconv.Atoi(selectedVersionInput)
	selectedVersion := installedVersions[selectedVersionIndex]

	space()

	// construct the path to the target file and validate that it exists
	appAsar := filepath.FromSlash(fmt.Sprintf("%s/%s/resources/app.asar", path, selectedVersion))

	fmt.Printf("Locating patch target file %s...", appAsar)

	if _, err := os.Stat(appAsar); os.IsNotExist(err) {
		panic(fmt.Sprintf("Unpatchable Slack installation.  Slack may have changed things, or your install is broken."))
	}

	fmt.Print(" done.\n")

	space()

	// fetch the code to use for the patch
	if *srcURLInputPtr != "" {
		srcURL = *srcURLInputPtr
		fmt.Print(fmt.Sprintf("Fetching source from %s...", srcURL))
	} else {
		fmt.Print(fmt.Sprintf("Fetching source from default location %s...", srcURL))
	}

	getPatchSrc(srcURL)

	fmt.Print(" done.\n")

	space()

	// extract the file
	fmt.Printf("Extracting target archive...")

	// locate .\dist\ssb-interop.bundle.js

	// read the final line in the file and check whether it matches '//# sourceMappingURL=ssb-interop.bundle.js.map'
	// if so, the file has already been patched.  bail out after deleting the extracted files.

	// append the fetched patch code to the file

	// create a new archive

	// rename the old file

	// rename the new archive to the old filename

	// done!
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

func getPatchSrc(srcURL string) (src string) {
	resp, err := http.Get(srcURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch url %s: %s", srcURL, err.Error))
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}
