package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gplus"
)

type templateHandler struct {
	once     sync.Once
	fileName string
	template *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.template = template.Must(template.ParseFiles(filepath.Join("templates", t.fileName)))
	})
	t.template.Execute(w, r)
}

var (
	googleClientID    = os.Getenv("GOOGLE_CLIENT_ID")
	googleSecret      = os.Getenv("GOOGLE_SECRET")
	googleRedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")
)

func main() {

	log.Println("Start server")
	var addr = flag.String("addr", ":8080", "address of application")
	flag.Parse()

	// oauth setup
	goth.UseProviders(
		gplus.New(googleClientID, googleSecret, googleRedirectURL),
	)

	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHandler{fileName: "chat.html"}))
	http.Handle("/login", &templateHandler{fileName: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// get the room going
	go r.run()
	// start the web server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err.Error())
	}
}
