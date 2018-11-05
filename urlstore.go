package main

import (
	"encoding/json"
	"io"
)

type urlStore struct {
	reader  io.Reader
	urlsets [][]string
}

var pathToUrls = "/home/steve/craigslisturls.json"

func (s *urlStore) loadURLs() error {
	n, err := external.readfile(pathToUrls)
	if err != nil {
		return err
	}

	err = s.parseURLsFile(n)
	return err
}

func (s *urlStore) parseURLsFile(b []byte) error {

	err := json.Unmarshal(b, &s.urlsets)
	return err
}

func (s *urlStore) setURLAt(setIndex, urlIndex int, url string) {
	s.urlsets[setIndex][urlIndex] = url
	b, _ := json.Marshal(s.urlsets)
	external.writefile(string(b))
}

func (s *urlStore) deleteURLAt(index int) {
}
func (s *urlStore) addURL() {
	//append string slice with dummy string value
}

func (s *urlStore) getUrls(setIndex int) []string {
	return nil
}

func (s *urlStore) touch(filename string) {
}
