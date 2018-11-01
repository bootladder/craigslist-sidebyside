package main

import (
	"errors"
	"testing"
)

func Test_loadURLs_readFailed_returnsError(t *testing.T) {
	var urlstore urlStore
	mockreader := mockReader{}
	mockreader.shouldError = true
	urlstore.reader = mockreader

	err := urlstore.loadURLs()
	if err == nil {
		t.Errorf("expected file io error, got nil")
	}
}
func Test_loadURLs_badJSON_returnsError(t *testing.T) {
	var urlstore urlStore
	mockreader := mockReader{}
	mockreader.shouldError = false
	mockreader.returnedBytes = []byte("hello")
	urlstore.reader = mockreader

	err := urlstore.loadURLs()
	if err == nil {
		t.Errorf("expected file io error, got nil")
	}
}

func Test_loadURLs_goodJSON_OK(t *testing.T) {

	jsonString := `
	{
		"arra": [
			"http://boston.craigslist.org/jjj/?query=hello",
			"http://portland.craigslist.org/jjj/?query=hello"
		],
		"arr": [
			"http://austin.craigslist.org/jjj/?query=hello",
			"http://houston.craigslist.org/jjj/?query=hello"
		]
	}
		`
	var urlstore urlStore
	mockreader := mockReader{}
	mockreader.shouldError = false
	mockreader.returnedBytes = []byte(jsonString)
	urlstore.reader = mockreader

	err := urlstore.loadURLs()
	if err != nil {
		t.Errorf("expected no error, got error: " + err.Error())
	}

	//expected := urlstore.urlsets[0].urls[0]
	//if expected != "http://boston.craigslist.org/jjj/?query=hello" {
	//	t.Errorf("expected string got no string")
	//}
}
func Test_loadURLs_topJSONIsNotArray_Fails(t *testing.T) {

	jsonString := `
	{
		"validjson" : "butNotAnArray"
	}
	
	`
	var urlstore urlStore
	mockreader := mockReader{}
	mockreader.shouldError = false
	mockreader.returnedBytes = []byte(jsonString)
	urlstore.reader = mockreader

	err := urlstore.loadURLs()
	if err != nil {
		t.Errorf("expected no error, got error")
	}

}

type mockReader struct {
	shouldError   bool
	returnedBytes []byte
}

func (m mockReader) Read(p []byte) (n int, err error) {
	if m.shouldError == true {
		return 0, errors.New("hello")
	}
	copy(p, m.returnedBytes)
	return len(p), nil
}
