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
	s.saveUrlSets()
}

func (s *urlStore) saveUrlSets() {
	b, _ := json.Marshal(s.urlsets)
	external.writefile(string(b))
}

func (s *urlStore) deleteURLAt(setIndex, urlIndex int) {
	set := s.urlsets[setIndex]
	set = append(set[:urlIndex], set[(urlIndex+1):]...)
	s.urlsets[setIndex] = set
	s.saveUrlSets()
}
func (s *urlStore) addURL(setIndex int) {
	s.urlsets[setIndex] = append(s.urlsets[setIndex], "http://boston.craigslist.org/jjj/?query=hello")
}

func (s *urlStore) getUrls(setIndex int) []string {
	return nil
}

func (s *urlStore) touch(filename string) {
}
