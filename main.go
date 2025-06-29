package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jcoene/go-base62"
)

var m = make(map[string]string)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := mux.NewRouter()
	r.HandleFunc("/shorten/{long_link}", ShortenLink)
	r.HandleFunc("/favicon.ico", GetIcon)
	r.HandleFunc("/{short_link}", Transfer)
	http.Handle("/", r)
	log.Println("Сервер запущен на :9090")
	http.ListenAndServe(":9090", nil)
}

func GetIcon(w http.ResponseWriter, r *http.Request) {}

func ShortenLink(w http.ResponseWriter, r *http.Request) {
	var long_link = mux.Vars(r)["long_link"]
	log.Printf("Ссылка для сокращения %s", long_link)

	if !strings.HasPrefix(long_link, "http://") && !strings.HasPrefix(long_link, "https://") {
		long_link = "http://" + long_link
	}

	m[base62.Encode(int64(len(m)))] = long_link
	fmt.Fprintf(w, "Long link is %v\n", long_link)
	fmt.Fprintf(w, "Current map is %v\n", m)
}

func Transfer(w http.ResponseWriter, r *http.Request) {
	var short_link = mux.Vars(r)["short_link"]

	long_link, exists := m[short_link]
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		log.Println(short_link, " not found in ", m)
		return
	}

	http.Redirect(w, r, long_link, http.StatusFound)

}
