package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"urlshortener/app/repo/link"

	"github.com/gorilla/mux"
)

type Handler struct {
	*mux.Router
	l *link.Links
}

var Templates struct {
	Main *template.Template
	Msg  *template.Template
}

type Msg struct {
	Msg string
}

func NewHandlers(l *link.Links) *Handler {
	var err error
	Templates.Main, err = template.ParseFiles("api/templates/main.html")
	if err != nil {
		log.Fatalf("Init template error: %v", err)
	}
	Templates.Msg, err = template.ParseFiles("api/templates/msg.html")
	if err != nil {
		log.Fatalf("Init template error: %v", err)
	}

	router := mux.NewRouter()
	ret := &Handler{
		router,
		l,
	}
	router.HandleFunc("/", ret.MainPage).Methods(http.MethodGet)
	router.HandleFunc("/", ret.CreateLink).Methods(http.MethodPost)
	router.HandleFunc("/{link}", ret.ReadLink).Methods(http.MethodGet)

	return ret
}

func (rt *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	err := Templates.Main.ExecuteTemplate(w, "main", struct{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rt *Handler) CreateLink(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	longLink := r.PostFormValue("long_link")

	// check that the long link is not empty before create short link
	if longLink == "" {
		err := Templates.Msg.ExecuteTemplate(w, "msg", Msg{"link must be not empty"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	shortLink, err := rt.l.Create(r.Context(), longLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Templates.Msg.ExecuteTemplate(w, "msg", Msg{fmt.Sprint("http://localhost:8000/", shortLink)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rt *Handler) ReadLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortLink := vars["link"]

	// check that the short link is not empty before search long link
	if shortLink == "" {
		http.Error(w, "empty link", http.StatusBadRequest)
		return
	}

	longLink, err := rt.l.Read(r.Context(), shortLink)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "read link error", http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, longLink, http.StatusFound)
}
