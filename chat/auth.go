package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/stretchr/objx"
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

// loginHandler handles the third-party login process.
// format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	query := r.URL.Query()
	query.Set("provider", provider)
	r.URL.RawQuery = query.Encode()
	switch action {
	case "login":
		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
			// do some thing after login
			setCookie(gothUser, w)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
	case "callback":
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		setCookie(user, w)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func setCookie(userInfo goth.User, w http.ResponseWriter) {
	m := md5.New()
	io.WriteString(m, strings.ToLower(userInfo.Email))
	userID := fmt.Sprintf("%x", m.Sum(nil))
	authCookieValue := objx.New(map[string]interface{}{
		"user_id":    userID,
		"name":       userInfo.Name,
		"avatar_url": userInfo.AvatarURL,
		"email":      userInfo.Email,
	}).MustBase64()
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/",
	})
	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// MustAuth is a Auth hanlder wrapper
func MustAuth(h http.Handler) http.Handler {
	return &authHandler{next: h}
}
