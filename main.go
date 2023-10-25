package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"

	"github.com/msa-ali/picbucket/controllers"
	"github.com/msa-ali/picbucket/migrations"
	"github.com/msa-ali/picbucket/models"
	"github.com/msa-ali/picbucket/templates"
	"github.com/msa-ali/picbucket/utils"
	"github.com/msa-ali/picbucket/views"
)

func main() {
	// LOAD ENV
	err := utils.LoadEnv()
	if err != nil {
		panic(err)
	}
	env := utils.GetEnv()

	// Setup database and do migration
	db, err := models.Open(models.DefaultPostgresConfig())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Setup services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	galleryService := &models.GalleryService{
		DB: db,
	}
	emailService, err := models.NewEmailService(models.SMTPConfig{
		Host:     env.SMTPHost,
		Port:     env.SMTPPort,
		Username: env.SMTPUsername,
		Password: env.SMTPPassword,
	})
	if err != nil {
		panic(err)
	}

	// Setup middlewares
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}
	csrfMiddleware := csrf.Protect(
		[]byte(env.CSRFKey),
		csrf.Secure(env.CSRFSecure),
		csrf.Path("/"),
	)

	// Setup Controller
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}
	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "forgot-password.gohtml", "tailwind.gohtml"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "check-your-email.gohtml", "tailwind.gohtml"))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "reset-pw.gohtml", "tailwind.gohtml"))
	galleriesC.Templates.Index = views.Must(
		views.ParseFS(
			templates.FS,
			path.Join("./galleries", "index.gohtml"),
			"tailwind.gohtml",
		),
	)
	galleriesC.Templates.New = views.Must(
		views.ParseFS(
			templates.FS,
			path.Join("./galleries", "new.gohtml"),
			"tailwind.gohtml",
		),
	)
	galleriesC.Templates.Edit = views.Must(
		views.ParseFS(
			templates.FS,
			path.Join("./galleries", "edit.gohtml"),
			"tailwind.gohtml",
		),
	)
	galleriesC.Templates.Show = views.Must(
		views.ParseFS(
			templates.FS,
			path.Join("./galleries", "show.gohtml"),
			"tailwind.gohtml",
		),
	)

	// setup routes
	r := chi.NewRouter()
	r.Use(csrfMiddleware)

	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "notfound.gohtml"))
	r.NotFound(controllers.StaticHandler(tpl))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.Authenticate)
	r.Post("/signout", usersC.SignOut)
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)
	// r.Get("/")
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleriesC.Show)
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/", galleriesC.Index)
			r.Get("/new", galleriesC.New)
			r.Post("/", galleriesC.Create)
			r.Get("/{id}/edit", galleriesC.Edit)
			r.Post("/{id}", galleriesC.Update)
		})
	})

	r.Get("/users/me", usersC.CurrentUser)
	port := utils.GetEnv().ServerPort
	fmt.Printf("Starting the server at port :%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), umw.SetUser(r))
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
