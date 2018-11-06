package main

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadURLs_readFailed_returnsError(t *testing.T) {
	var urlstore urlStore

	usingMockReadFileFail(errors.New("fail reading"))
	err := urlstore.loadURLs()
	assert.Error(t, err, "expected file io error, got nil")
}
func Test_loadURLs_badJSON_returnsError(t *testing.T) {
	var urlstore urlStore

	jsonString := `
	thisisnotjson
		`
	usingMockReadFileSuccess([]byte(jsonString))
	err := urlstore.loadURLs()
	assert.Error(t, err, "expected file io error, got nil")
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
	assert.Error(t, err, "expected file io error, got nil")
}

func Test_loadURLs_goodJSON_OK(t *testing.T) {

	jsonString := `
	[
		[
			"http://boston.craigslist.org/jjj/?query=hello",
			"http://portland.craigslist.org/jjj/?query=hello"
		],
		[
			"http://austin.craigslist.org/jjj/?query=hello",
			"http://houston.craigslist.org/jjj/?query=hello"
		]
	] `
	usingMockReadFileSuccess([]byte(jsonString))
	var urlstore urlStore

	err := urlstore.loadURLs()
	assert.NoError(t, err, "expected no error, got error: ")

	assert.Equal(t, urlstore.urlsets[0][0], "http://boston.craigslist.org/jjj/?query=hello")
	assert.Equal(t, urlstore.urlsets[0][1], "http://portland.craigslist.org/jjj/?query=hello")
	assert.Equal(t, urlstore.urlsets[1][0], "http://austin.craigslist.org/jjj/?query=hello")
	assert.Equal(t, urlstore.urlsets[1][1], "http://houston.craigslist.org/jjj/?query=hello")
}

func Test_setUrlAt_storesUrlInArray_andCallsSave(t *testing.T) {
	external.writefile = mockWriteFile
	var urlstore urlStore

	urlstore.urlsets = [][]string{{"orig1", "orig2"}, {"orig1", "orig2"}}
	expectedURLSets := [][]string{{"orig1", "orig2"}, {"orig1", "newurl"}}
	jsonBytes, _ := json.Marshal(expectedURLSets)

	urlstore.setURLAt(1, 1, "newurl")

	assert.Equal(t, urlstore.urlsets[1][1], "newurl")
	assert.Equal(t, mockwritefileBytes, jsonBytes)
}

func Test_deleteUrlAt_removesFromArray_andCallsSave(t *testing.T) {
	external.writefile = mockWriteFile
	var urlstore urlStore

	urlstore.urlsets = [][]string{{"orig1", "orig2"}, {"orig1", "orig2"}}
	expectedURLSets := [][]string{{"orig2"}, {"orig1", "orig2"}}
	jsonBytes, _ := json.Marshal(expectedURLSets)

	urlstore.deleteURLAt(0, 0)

	assert.Equal(t, urlstore.urlsets, expectedURLSets)
	assert.Equal(t, mockwritefileBytes, jsonBytes)
}

func Test_addURL_addsToArray_andCallsSave(t *testing.T) {
	external.writefile = mockWriteFile
	var urlstore urlStore

	urlstore.urlsets = [][]string{{"orig1"}, {"orig1", "orig2"}}
	expectedURLSets := [][]string{{"orig1", "http://boston.craigslist.org/jjj/?query=hello"}, {"orig1", "orig2"}}

	jsonBytes, _ := json.Marshal(expectedURLSets)

	urlstore.addURL(0)

	assert.Equal(t, expectedURLSets, urlstore.urlsets)
	assert.Equal(t, mockwritefileBytes, jsonBytes)
}

/////////////////////////////////////////////
var mockreadfileBytes []byte
var mockreadfileError error

var mockwritefileBytes []byte

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

func mockWriteFile(filename string, content []byte, mode os.FileMode) error {

	mockwritefileBytes = content
	return errors.New("hello")
}
