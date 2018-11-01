package main

import (
	"encoding/json"
	"io"
)

type urlStore struct {
	reader  io.Reader
	urlsets []urlSet
}
type urlSet struct {
	urls []string
}

func (s *urlStore) loadURLs() error {
	var b []byte
	n, err := s.reader.Read(b)
	printf("the url string is : %v and number of bytes read is %v", b, n)
	if err != nil {
		return err
	}

	err = s.parseURLsFile(b)
	return err
}

func (s *urlStore) parseURLsFile(b []byte) error {

	urlsets := make([]urlSet, 0)
	err := json.Unmarshal(b, urlsets)
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
