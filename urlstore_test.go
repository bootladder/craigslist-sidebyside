package main

import (
	"errors"
	"testing"
)

func Test_loadURLs_readFailed_returnsError(t *testing.T) {
	var urlstore urlStore

	usingMockReadFileFail(errors.New("fail reading"))
	err := urlstore.loadURLs()
	if err == nil {
		t.Errorf("expected file io error, got nil")
	}
}
func Test_loadURLs_badJSON_returnsError(t *testing.T) {
	var urlstore urlStore

	jsonString := `
	thisisnotjson
		`
	usingMockReadFileSuccess([]byte(jsonString))
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
	usingMockReadFileSuccess([]byte(jsonString))
	var urlstore urlStore

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

	usingMockReadFileSuccess([]byte(jsonString))

	err := urlstore.loadURLs()
	if err != nil {
		t.Errorf("expected no error, got error")
	}
}

var mockreadfileBytes []byte
var mockreadfileError error

func mockReadFile(filename string) ([]byte, error) {
	return mockreadfileBytes, mockreadfileError
}

func usingMockReadFileSuccess(myBytes []byte) {

	external.readfile = mockReadFile
	mockreadfileError = nil
	mockreadfileBytes = myBytes
}

func usingMockReadFileFail(err error) {

	external.readfile = mockReadFile
	mockreadfileError = err
	mockreadfileBytes = []byte("doesn't matter")
}
