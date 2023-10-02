package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/msa-ali/picbucket/controllers"
	"github.com/msa-ali/picbucket/views"
)

// func notFoundhandler(w http.ResponseWriter, r *http.Request) {
// 	// 1
// 	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	// w.WriteHeader(http.StatusNotFound)
// 	// fmt.Fprint(w, "<h1>404 - page not found</h1>")
// 	// 2
// 	// http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 	// 3
// 	tplPath := filepath.Join("templates", "notfound.gohtml")
// 	executeTemplate(w, tplPath)
// }

func main() {
	r := chi.NewRouter()
	tpl, err := views.Parse(filepath.Join("templates", "home.gohtml"))
	if err != nil {
		panic(err)
	}
	r.Get("/", controllers.StaticHandler(tpl))

	tpl, err = views.Parse(filepath.Join("templates", "contact.gohtml"))
	if err != nil {
		panic(err)
	}
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl, err = views.Parse(filepath.Join("templates", "faq.gohtml"))
	if err != nil {
		panic(err)
	}
	r.Get("/faq", controllers.StaticHandler(tpl))

	tpl, err = views.Parse(filepath.Join("templates", "notfound.gohtml"))
	if err != nil {
		panic(err)
	}
	r.NotFound(controllers.StaticHandler(tpl))
	fmt.Println("Starting the server at port :8080")
	http.ListenAndServe(":8080", r)
}
