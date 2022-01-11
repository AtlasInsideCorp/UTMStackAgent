package main

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func multiDownload(files []string, server string) error {
	for _, f := range files {
		err := download(server + f)
		if err != nil {
			return err
		}
	}
	return nil
}

func download(f string) error {
	var fileName string

	fileURL, err := url.Parse(f)
	if err != nil {
		return err
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName = segments[len(segments)-1]

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	// Put content on file
	resp, err := client.Get(f)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	defer file.Close()

	return err
}
