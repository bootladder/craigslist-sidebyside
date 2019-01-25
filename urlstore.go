package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type urlStore struct {
	reader  io.Reader
	urlsets []urlSet
}

type urlSet struct {
	Name string   `json:"name"`
	Urls []string `json:"urls"`
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
	s.urlsets[setIndex].Urls[urlIndex] = url
	s.saveURLSetsToFile()
}

func (s *urlStore) saveURLSetsToFile() {
	b, _ := json.MarshalIndent(s.urlsets, "", "  ")
	external.writefile(pathToUrls, b, 0)
}

func (s *urlStore) deleteURLAt(setIndex, urlIndex int) {
	urls := s.urlsets[setIndex].Urls
	urls = append(urls[:urlIndex], urls[(urlIndex+1):]...)
	s.urlsets[setIndex].Urls = urls
	s.saveURLSetsToFile()
}
func (s *urlStore) addURL(setIndex int) {
	s.urlsets[setIndex].Urls = append(s.urlsets[setIndex].Urls, "http://boston.craigslist.org/search/jjj/?sort=date&query=engineer")
	s.saveURLSetsToFile()
}

func (s *urlStore) addNewURLSet() {
	s.urlsets = append(s.urlsets, urlSet{})
}

func (s *urlStore) getAllURLSetNames() []string {

	var names = []string{}
	for _, urlset := range s.urlsets {
		names = append(names, urlset.Name)
	}
	return names
}

func (s *urlStore) touch(filename string) {
}
