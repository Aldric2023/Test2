package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type UserData struct {
	PageTitle string
	Body      template.HTML
}

func middlewareA(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//favicon.ico is a reuest made by the browser automatically, i will ignore it
		//just return when i encounter this request just to keep the log clean
		if r.URL.Path == "/favicon.ico" {
			return
		}
		log.Println("Executing middle ware A")

		next.ServeHTTP(w, r) //everything in this function is executed before the middleware function ends

		//this is executed on the way up
		log.Println("Executing middleware A again")

	})
}

func middlewareB(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//this is executed on the way down to the handeler
		log.Println("Executing middleware B")

		log.Printf("IP address: %s ", r.RemoteAddr)
		log.Println("URL Entered: ", r.URL.Path)

		if r.URL.Path != "/" {
			//return and dont go any further
			fmt.Println("invalid link")
			return
		}

		//we now run our homehandeler, the function we passed into the middleware

		next.ServeHTTP(w, r) //everything in this function is executed before the middleware function ends

		//this is executed on the way up
		log.Println("Executing middleware B again")

	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	body := "<h2><p>learning about middle ware<br>" +
		"This is the main home page!!! If you got here then that means you used the right link<br>" +
		"</p></h2>"

	data := UserData{
		PageTitle: "Test Page",
		Body:      template.HTML(body),
	}

	ts, _ := template.ParseFiles("public/index.html.tmpl")

	ts.Execute(w, data)

	log.Println("home page working")

}

func main() {
	// serve static files from the "static" directory
	mux := http.NewServeMux()

	// file server
	fs := http.FileServer(http.Dir("public"))

	// to access images from the public dir
	mux.Handle("/public/", http.StripPrefix("/public/", fs))

	//handeler that serves the home page
	mux.Handle("/", middlewareA(middlewareB(http.HandlerFunc(homeHandler))))

	// start the server on port 8080 on localhost
	log.Fatal(http.ListenAndServe("localhost:8080", mux))

}
