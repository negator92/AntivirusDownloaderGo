package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	fileName    string
	fullURLFile string
)

func main() {

	currentTime := time.Now()
	fullURLFiles := []string{
		"https://rescuedisk.s.kaspersky-labs.com/latest/krd.iso",
		"https://download.geo.drweb.com/pub/drweb/livedisk/drweb-livedisk-900-cd.iso",
	}

	for _, fullURLFile = range fullURLFiles {

		// Build fileName from fullPath
		fileURL, err := url.Parse(fullURLFile)
		if err != nil {
			log.Fatal(err)
		}
		path := fileURL.Path
		segments := strings.Split(path, "/")
		fileName = segments[len(segments)-1]

		//drweb contains "-", cut it
		if strings.Contains(fileName, "-") {
			fileNamePreParts := strings.Split(fileName, "-")
			fileName = fileNamePreParts[0] + ".iso"
		}

		//paste date into fileName
		fileNameParts := strings.Split(fileName, ".")
		fileNameParts = append(fileNameParts, "-"+currentTime.Format("2006-01-02")+".")
		fileNameParts[len(fileNameParts)-2], fileNameParts[len(fileNameParts)-1] = fileNameParts[len(fileNameParts)-1], fileNameParts[len(fileNameParts)-2]
		fileName = fileNameParts[0] + fileNameParts[len(fileNameParts)-2] + fileNameParts[len(fileNameParts)-1]
		// Create blank file
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		client := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
		// Put content on file
		resp, err := client.Get(fullURLFile)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		size, err := io.Copy(file, resp.Body)

		defer file.Close()

		fmt.Printf("Downloaded a file %s with size %d\n", fileName, size)
	}

}
