package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
)

func main() {
	path := getSlackDir()

	fmt.Println(fmt.Sprintf("Slack directory: %s", path))

	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic("error reading Slack directory.")
	}

	for _, f := range files {
		name := f.Name()
		isAppDir, _ := regexp.MatchString("^app-((\\d+\\.?){3})$", name)

		if f.IsDir() && isAppDir {
			fmt.Println(fmt.Sprintf("Detected Slack version: %s", name))
		}
	}
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
