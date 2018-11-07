package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
)

var debug = true

var err error

var urlstore urlStore

type craigslistRequest struct {
	SearchURL   string `json:"searchURL"`
	ColumnIndex int    `json:"columnIndex"`
	SetIndex    int    `json:"setIndex"`
}
type craigslistResponse struct {
	ResponseHTML string `json:"response"`
}

type craigslistDeleteRequest struct {
	ColumnIndex int `json:"columnIndex"`
	SetIndex    int `json:"setIndex"`
}

type craigslistAddRequest struct {
	SetIndex int `json:"setIndex"`
}

type craigslistGetRequest struct {
	SetIndex int `json:"setIndex"`
}
type returnURLSetResponse struct {
	Urls []string `json:"urls"`
}

func main() {

	inject()
	urlstore.loadURLs()

	router := httprouter.New()
	router.ServeFiles("/static/*filepath",
		http.Dir("public"))

	router.POST("/api/", createPostHandler(""))
	router.GET("/api/:setIndex", getURLSet)
	router.DELETE("/api/", createDeleteHandler(""))
	router.PUT("/api/", createPutHandler(""))

	browser.OpenURL("http://localhost:8080/static/index.html")
	http.ListenAndServe(":8080", router)
}

func postNoteHandler(w http.ResponseWriter, r *http.Request) {

	var req craigslistRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	fatal(err)
	req.SearchURL, err = url.QueryUnescape(req.SearchURL)
	fatal(err)
	fmt.Printf("POST URL: index is %d\n", req.ColumnIndex)

	var resp craigslistResponse

	if debug == true {
		resp.ResponseHTML = "<b>HELLO RESPONSE</b>"
	} else {
		resp.ResponseHTML = makeRequest(req.SearchURL)
	}

	urlstore.setURLAt(req.SetIndex, req.ColumnIndex, req.SearchURL)

	jsonOut, err := json.Marshal(resp)
	fatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonOut)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {

	var req craigslistDeleteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	fatal(err)
	fmt.Printf("Delete: the index is %d\n", req.ColumnIndex)

	urlstore.deleteURLAt(req.SetIndex, req.ColumnIndex)

	returnURLSetJSONResponse(w, req.SetIndex)
}

func putHandler(w http.ResponseWriter, r *http.Request) {

	var req craigslistAddRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	fatal(err, "JSON Decode Put Handler Body")
	fmt.Printf("Add(Put): the index is %d\n", req.SetIndex)
	urlstore.addURL(req.SetIndex)
	returnURLSetJSONResponse(w, req.SetIndex)
}

func getURLSet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	setIndexString := ps.ByName("setIndex")
	setIndex, err := strconv.Atoi(setIndexString)
	fatal(err)
	returnURLSetJSONResponse(w, setIndex)
}

func makeRequest(url string) string {
	log.Println("makeRequest: " + url)
	resp, err := http.Get(url) //"https://httpbin.org/get"

	//gracefully handle error with invalid craigslist URL
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}

func returnURLSetJSONResponse(w http.ResponseWriter, setIndex int) {
	var resp returnURLSetResponse
	resp.Urls = urlstore.urlsets[setIndex]

	jsonOut, err := json.MarshalIndent(resp, "", "  ")
	fatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonOut)
}

func fatal(err error, msgs ...string) {
	if err != nil {
		var str string
		for _, msg := range msgs {
			str = msg
			break
		}
		panic(errors.Wrap(err, str))
	}
}

func printf(s string, a ...interface{}) {
	fmt.Printf(s, a...)
}

func createPostHandler(msg string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		postNoteHandler(w, r)
	}
}

func createDeleteHandler(msg string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		deleteHandler(w, r)
	}
}

func createPutHandler(msg string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		putHandler(w, r)
	}
}

var external = externalFuncs{}

type externalFuncs struct {
	readfile  func(string) ([]byte, error)
	writefile func(string, []byte, os.FileMode) error
}

func inject() {
	external.readfile = ioutil.ReadFile
	external.writefile = ioutil.WriteFile
}
