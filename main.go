package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

// Bug doc
type Bug struct {
	Name        string        `json:"name"`
	URL         string        `json:"url"`
	Description string        `json:"description"`
	Keywords    []interface{} `json:"keywords"`
	Code        struct {
		Language string `json:"language"`
		URL      string `json:"url"`
		Files    []struct {
			FileName string `json:"file_name"`
			Patch    string `json:"patch"`
		} `json:"files"`
	} `json:"code"`
}

var delimiter = strings.Repeat("/", 128)

func main() {
	// Create the results directory
	err := os.MkdirAll("results", 0775)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	files, err := ioutil.ReadDir("./data")
	for _, fileInfo := range files {
		// Open the repo list file
		data, err := ioutil.ReadFile("./data/" + fileInfo.Name())
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// Load the list of bugs from the json file
		var bugs []Bug
		json.Unmarshal(data, &bugs)
		fileOut, err := os.Create("./results/" + fileInfo.Name() + "_patches.txt")

		for _, bug := range bugs {
			for _, file := range bug.Code.Files {
				if err != nil {
					log.Fatal("Cannot create file", err)
				}
				defer fileOut.Close()
				lines := strings.Split(file.Patch, "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "@@") {
						var re = regexp.MustCompile(`^@@\s.+?\s@@`)
						line = re.ReplaceAllString(line, ``)
						if line == "" {
							continue
						}
					}
					if !strings.HasPrefix(line, "+") {
						if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "-") {
							line = line[1:]
						}
						if !strings.HasPrefix(line, `\ No newline at end of file`) {
							fmt.Fprintf(fileOut, line+"\n")
						}
					}
				}
			}
			fmt.Fprintf(fileOut, delimiter+"\n")
		}
	}

}
