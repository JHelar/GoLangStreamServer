package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//Create person
type person struct {
	fName string
}

//Attach function to person struct
func (p *person) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	log.Println(path)
	data, err := ioutil.ReadFile(string(path))
	if err == nil {
		var contentType string
		if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".html") {
			contentType = "text/html"
		} else if strings.HasSuffix(path, ".js") {
			contentType = "application/javascript"
		} else if strings.HasSuffix(path, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(path, ".svg") {
			contentType = "image/svg+xml"
		} else {
			contentType = "text/plain"
		}
		w.Header().Add("Content Type", contentType)
		w.Write(data)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 " + p.fName + " My Friend - " + http.StatusText(404)))
	}
}

func main() {
	personOne := &person{fName: "John"}
	http.Handle("/", personOne)
	http.ListenAndServe(":8080", personOne)
}
 