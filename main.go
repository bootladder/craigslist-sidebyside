package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
)

var err error

type craigslistRequest struct {
	SearchURL string `json:"searchURL"`
}

type craigslistResponse struct {
	ResponseHTML string `json:"response"`
}

type getUrlsResponse struct {
	Urls []string `json:"urls"`
}

func main() {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath",
		http.Dir("public"))

	router.POST("/api/", createPostHandler(""))
	router.GET("/api/", createGetHandler(""))

	browser.OpenURL("http://localhost:8080/static/index.html")
	http.ListenAndServe(":8080", router)
}

func postNoteHandler(w http.ResponseWriter, r *http.Request) {

	var req craigslistRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	req.SearchURL, _ = url.QueryUnescape(req.SearchURL)

	var resp craigslistResponse
	resp.ResponseHTML = makeRequest(req.SearchURL)

	//Save the URL
	saveURL(req.SearchURL)

	jsonOut, err := json.Marshal(resp)
	fatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonOut)
}

func getHandler(w http.ResponseWriter, r *http.Request) {

	var resp getUrlsResponse

	//Get URLS from persistent storage
	resp.Urls = loadURLs()

	jsonOut, err := json.Marshal(resp)
	fatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonOut)
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

	//log.Println(string(body))
	return string(body)
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
	fmt.Printf(s, a)
}

func createPostHandler(msg string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		postNoteHandler(w, r)
	}
}

func createGetHandler(msg string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		getHandler(w, r)
	}
}
