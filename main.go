package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func executeTemplate(w http.ResponseWriter, filePath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	template, err := template.ParseFiles(filePath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
		return
	}
	err = template.Execute(w, nil)
	if err != nil {
		log.Printf("error while executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
}

func homehandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tplPath)
}

func notFoundhandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// w.WriteHeader(http.StatusNotFound)
	// fmt.Fprint(w, "<h1>404 - page not found</h1>")
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tplPath)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "faq.gohtml")
	executeTemplate(w, tplPath)
}

func main() {
	r := chi.NewRouter()
	r.Get("/", homehandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(notFoundhandler)
	fmt.Println("Starting the server at port :8080")
	http.ListenAndServe(":8080", r)
}
