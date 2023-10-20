package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"

	"github.com/msa-ali/picbucket/controllers"
	"github.com/msa-ali/picbucket/migrations"
	"github.com/msa-ali/picbucket/models"
	"github.com/msa-ali/picbucket/templates"
	"github.com/msa-ali/picbucket/views"
)

func main() {
	// init Database
	db, err := models.Open(models.DefaultPostgresConfig())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	// user controller
	usersC := controllers.Users{
		UserService: &models.UserService{
			DB: db,
		},
		SessionService: &models.SessionService{
			DB: db,
		},
	}
	r := chi.NewRouter()

	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)

	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.Authenticate)
	r.Post("/signout", usersC.SignOut)
	r.Get("/users/me", usersC.CurrentUser)

	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "notfound.gohtml"))
	r.NotFound(controllers.StaticHandler(tpl))

	fmt.Println("Starting the server at port :8080")
	csrfKey := "dmUQrNkHKnGBrdovbeLNNqjAIzinTVDa"
	csrfMiddleware := csrf.Protect([]byte(csrfKey), csrf.Secure(false))
	http.ListenAndServe(":8080", csrfMiddleware(r))
}

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

// func TimerMiddleware(h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		h(w, r)
// 		fmt.Printf("Request time for %s - %v\n\n", r.URL.Path, time.Since(start))
// 	}
// }
