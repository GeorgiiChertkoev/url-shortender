package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcoene/go-base62"
)

var m = make(map[string]string)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := mux.NewRouter()
	r.HandleFunc("/shorten", ShortenLink)
	r.HandleFunc("/favicon.ico", GetIcon)
	r.HandleFunc("/{short_link}", Transfer)
	r.HandleFunc("/", IndexPage)
	http.Handle("/", r)
	log.Println("Сервер запущен на :9090")
	http.ListenAndServe(":9090", nil)
}

func GetIcon(w http.ResponseWriter, r *http.Request) {}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ShortenLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL не может быть пустым", http.StatusBadRequest)
		return
	}

	shortCode := base62.Encode(int64(len(m)))
	log.Printf("Ссылка для сокращения %s", originalURL)

	m[shortCode] = originalURL

	data := struct {
		OriginalURL string
		ShortURL    string
	}{
		OriginalURL: originalURL,
		ShortURL:    fmt.Sprintf("http://%s/%s", r.Host, shortCode),
	}

	tmpl, err := template.ParseFiles("templates/result.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)

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
