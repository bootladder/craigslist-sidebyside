package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type urlStore struct {
	reader  io.Reader
	urlsets [][]string
}

var pathToUrls = "/home/steve/craigslisturls.json"

func (s *urlStore) loadURLs() error {
	bytes, err := external.readfile(pathToUrls)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &s.urlsets)
	fmt.Printf("parseURLsFile: %v\n", s.urlsets)
	return err
}

func (s *urlStore) setURLAt(setIndex, urlIndex int, url string) {
	s.urlsets[setIndex][urlIndex] = url
	s.saveURLSetsToFile()
}

func (s *urlStore) saveURLSetsToFile() {
	b, _ := json.MarshalIndent(s.urlsets, "", "  ")
	external.writefile(pathToUrls, b, 0)
}

func (s *urlStore) deleteURLAt(setIndex, urlIndex int) {
	set := s.urlsets[setIndex]
	set = append(set[:urlIndex], set[(urlIndex+1):]...)
	s.urlsets[setIndex] = set
	s.saveURLSetsToFile()
}
func (s *urlStore) addURL(setIndex int) {
	s.urlsets[setIndex] = append(s.urlsets[setIndex], "http://boston.craigslist.org/search/jjj/?sort=date&query=engineer")
	s.saveURLSetsToFile()
}

func (s *urlStore) addNewURLSet() {
	s.urlsets = append(s.urlsets, make([]string, 0))
}

func (s *urlStore) getAllURLSetNames() []string {

	var names = []string{"hello", "hello2", "hello3"}
	return names
}

func (s *urlStore) touch(filename string) {
}
