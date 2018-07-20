package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
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

func main() {
	log.Println("Start server")
	var addr = flag.String("addr", ":8080", "address of application")
	flag.Parse()
	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHandler{fileName: "chat.html"}))
	http.Handle("/room", r)
	// get the room going
	go r.run()
	// start the web server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err.Error())
	}
}
