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
	"github.com/stretchr/objx"
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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.template.Execute(w, data)
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

	r := newRoom(UseGravatarAvatar)
	http.Handle("/chat", MustAuth(&templateHandler{fileName: "chat.html"}))
	http.Handle("/login", &templateHandler{fileName: "login.html"})
	http.Handle("/upload", &templateHandler{fileName: "upload.html"})
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	http.HandleFunc("/uploader", uploadHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// get the room going
	go r.run()
	// start the web server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err.Error())
	}
}
