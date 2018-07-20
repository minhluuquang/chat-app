package main

import (
	"net/http"
)

type authHandler struct {
	next http.Handler
}

func (a *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// success, call next handler
	a.next.ServeHTTP(w, r)
}

// MustAuth is a Auth hanlder wrapper
func MustAuth(h http.Handler) http.Handler {
	return &authHandler{next: h}
}
