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
	printf("the url string is : %v", string(n))
	if err != nil {
		return err
	}

	err = s.parseURLsFile(n)
	return err
}

func (s *urlStore) parseURLsFile(b []byte) error {

	var urlsets [][]string
	err := json.Unmarshal(b, &urlsets)
	urlstore.urlsets = urlsets
	return err
}

func (s *urlStore) save() {
	//json.marshal
	//writer.write()
}

func (s *urlStore) setURLAt(setIndex, urlIndex int, url string) {
	s.urlsets[setIndex][urlIndex] = url
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
