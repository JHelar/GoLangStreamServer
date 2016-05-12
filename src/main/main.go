package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"fmt"
	"encoding/hex"
	"time"
)

//Create person
type controller interface {
	getTemplate() string
}
type homeController struct {
}
type noteToHexController struct {
}
type streamController struct {
}
type streamImageController struct {
	isInit bool
	frameCount int
}

func (h *homeController) getTemplate() string {
	return "templates/home.html"
}
func (n *noteToHexController) getTemplate() string {
	return "templates/NoteToHex.html"
}
func (s *streamController) getTemplate() string {
	return "templates/Stream.html"
}


var intervalls = map[string]byte {
	"0":0x00,
	"1":0x40,
	"2":0x80,
	"3":0xC0,
}

var okt = map[string]byte {
	"0":0x00,
	"1":0x10,
	"2":0x20,
	"3":0x30,
}

var note = map[string]byte {
	"c":0x00,
	"c#":0x01,
	"d":0x02,
	"d#":0x03,
	"e":0x04,
	"f":0x05,
	"f#":0x06,
	"g":0x07,
	"g#":0x08,
	"a":0x09,
	"a#":0x0A,
	"b":0x0B,
}

func (n *noteToHexController) ParseNote(noteStr string) string{
	parts := strings.Split(noteStr, "-")
	result := [1]byte { 0 }

	if len(parts) == 3 {
		inter := intervalls[parts[0]]
		oktav := okt[parts[1]]
		not := note[parts[2]]
		result[0] = inter | oktav | not
	} 
	return hex.EncodeToString(result[:])
}

func (h *homeController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	//log.Println(path)
	if len(path) == 0 {
		path = h.getTemplate()
	} 
	handleRequest(w, r, path)
}

func (s *streamController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, s.getTemplate())
}

func (s *streamImageController) ServeHTTP(w http.ResponseWriter, r *http.Request){
	//Start go rutine to stream image to client.
	w.Header().Add("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	frameCount := 1
	for {
		//Get frame, write it to response
		//Using Motion JPG in order to stream new pictures to the page.
		file := fmt.Sprintf("public/img/GifTest/%v.jpeg", frameCount)
		//log.Println(file)
		frame, err := ioutil.ReadFile(file)
		if err == nil {
			frameCount = (frameCount % 5) + 1
			//Add boundery tag then contenttype and picture.
			_, wErr := fmt.Fprintf(w, "--frame\r\nContent-Type: image/jpeg\r\n\r\n")	
			if wErr != nil {
				log.Println("Cant write header to client, closing routine.")
				log.Println(wErr)
				break
			}else {
				//Write frame to client
				_, wErr := w.Write(frame)
				if wErr != nil {
					log.Println("Cant write frame to client, closing routine.")
					log.Println(wErr)
					break
				}else {
					time.Sleep(33 * time.Millisecond)
				}
			}
		}else{
			log.Println("Error in file read")
		}
	}
}
func (n *noteToHexController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		//w.Write(n.ParseNote(r.Form["note"][0]))
		fmt.Fprintf(w, "0x" + n.ParseNote(r.Form["note"][0]))
	} else {
		handleRequest(w, r, n.getTemplate())
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request, path string){
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
		log.Println(contentType)
		w.Header().Add("Content-Type", contentType)
		w.Write(data)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404  - " + http.StatusText(404)))
	}
}

//Attach function to person struct


func main() {
	homeController := new(homeController)
	noteToHexController := new(noteToHexController)
	streamController := new(streamController)
	streamImageController := new(streamImageController)

	http.Handle("/Stream/video_feed", streamImageController)
	http.Handle("/Stream", streamController)
	http.Handle("/NoteToHex", noteToHexController)	
	http.Handle("/", homeController)
	
	http.ListenAndServe(":8080", nil)
}
 