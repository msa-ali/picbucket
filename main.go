package main

import (
	"fmt"
	"net/http"
)

func homehandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to picbucket</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
		<h1>Contact Page</h1>
		<p>
			You can get in touch with me at 
				<a 
					href="mailto:altamashattari786@gmail.com"
				>
				altamashattari786@gmail.com
			</a>
		</p>
	`)
}

func main() {
	http.HandleFunc("/", homehandler)
	http.HandleFunc("/contact", contactHandler)
	fmt.Println("Starting the server at port :8080")
	http.ListenAndServe(":8080", nil)
}
