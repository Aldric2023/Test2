package main

import(
	"log"
	"fmt"
	"html/template"
	"net/http"
)

type UserData struct {
	PageTitle string
	Body      template.HTML
}


func middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//this is executed on the way down to the handeler
		log.Println("Executing middleware")
		log.Println(r.URL.Path)

		if r.URL.Path != "/" {
			//return and dont go any further
			fmt.Println("invalid link")
			return
		}

		next.ServeHTTP(w, r)
		//this is executed on the way up
		log.Println("Executing middleware again")
		log.Printf("IP address: %s ", r.RemoteAddr)
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

	// log.Println("home page working")

}


func main() {
	// serve static files from the "static" directory
	mux := http.NewServeMux()

	// file server
	fs := http.FileServer(http.Dir("public"))

	// to access images from the public dir
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	//handeler that serves the home page
	mux.Handle("/", middleware(http.HandlerFunc(homeHandler)))

	// start the server on port 8080 on localhost
	log.Fatal(http.ListenAndServe("localhost:8080", mux))

}
