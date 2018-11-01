package main

import (
	"encoding/json"
	"io"
)

type urlStore struct {
	reader  io.Reader
	urlsets urlSets
}
type urlSets struct {
	sets []urlSet
}
type urlSet struct {
	urls []string
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

	urlsets := urlSets{}
	err := json.Unmarshal(b, &urlsets)
	urlstore.urlsets = urlsets
	return err
	//return errors.New("hello")
}

func (s *urlStore) saveURL(url string) {
}

func (s *urlStore) save() {
}

func (s *urlStore) setURLAt(index int, url string) {
}
func (s *urlStore) deleteURLAt(index int) {
}
func (s *urlStore) addURL() {
}

func (s *urlStore) getUrls() []string {
	return nil
}

func (s *urlStore) loadURLSet2() {
}
func (s *urlStore) touch(filename string) {
}
