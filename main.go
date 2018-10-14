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

type note struct {
	Hello     string `json:"hello"`
	SearchURL string `json:"searchURL"`
	Response  string `json:"response"`
}

func main() {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath",
		http.Dir("public"))

	router.POST("/api/", createPostHandler(""))

	browser.OpenURL("http://localhost:8080/static/index.html")
	http.ListenAndServe(":8080", router)
}

func createPostHandler(msg string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		postNoteHandler(w, r)
	}
}

func postNoteHandler(w http.ResponseWriter, r *http.Request) {

	var note note
	err := json.NewDecoder(r.Body).Decode(&note)
	note.SearchURL, _ = url.QueryUnescape(note.SearchURL)

	note.Response = makeRequest(note.SearchURL)

	j, err := json.Marshal(note)
	fatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func makeRequest(url string) string {
	resp, err := http.Get(url) //"https://httpbin.org/get"
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
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
