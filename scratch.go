package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/browser"
)

func main4() {
	router := httprouter.New()
	// Mapping to methods is possible with HttpRouter
	router.GET("/api/", createGetHandler(""))
	router.POST("/api/", createPostHandler(""))
	// Path variable called name used here
	router.GET("/static/:name", getFileContent)

	fmt.Print("hello\n")
	browser.OpenURL("http://localhost:8080/static/index.html")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main1() {
	//r := mux.NewRouter().StrictSlash(false)
	//r.HandleFunc("/api/notes", GetNoteHandler).Methods("GET")
	//r.HandleFunc("/api/notes", PostNoteHandler).Methods("POST")

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)
	mux.HandleFunc("/api/notes", getNoteHandler)
	//mux.HandleFunc("/api/notes", PostNoteHandler)
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

	//http.ListenAndServe(":8080", r)
	http.ListenAndServe(":8080", mux)
}

func main2() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/", fs)
	mux.Handle("/welcome", &messageHandler{"hello123"})

	mux.Handle("/func", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "super lambda2")
		}))

	mux.Handle("/closure", createMessageHandlerFunc("my closure thing"))

	http.ListenAndServe(":8080", mux)
}

type messageHandler struct {
	message string
}

func (m *messageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, m.message)
}

func messageHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hellomessagefunc")
}

func createMessageHandlerFunc(msg string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, msg)
	})
}

// This is a function to execute a system command and return output
func getCommandOutput(command string, arguments ...string) string {
	// args... unpacks arguments array into elements
	cmd := exec.Command(command, arguments...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
	}
	return out.String()
}

func getFileContent(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, getCommandOutput("/bin/cat", "/home/steve/prog/go/src/scratch/stdlib-http-server/public/"+params.ByName("name")))
}

func createGetHandler(msg string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		getNoteHandler(w, r)
	}
}

func getNoteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("hello\n")
	var messages []string
	messages = append(messages, "arr1")
	messages = append(messages, "arr2")
	messages = append(messages, "arr3")
	messages = append(messages, "arr4")

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(messages)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
