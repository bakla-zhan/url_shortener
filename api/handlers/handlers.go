package handlers

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
	"urlshortener/app/repos/stat"
	"urlshortener/app/starter"

	"github.com/gorilla/mux"
)

type Handler struct {
	*mux.Router
	a *starter.App
}

var Templates struct {
	Main    *template.Template
	Msg     *template.Template
	Stats   *template.Template
	StatsIP *template.Template
}

type Msg struct {
	Msg1    string
	Msg2    string
	Msg3    string
	Msg4    string
	BackURL string
}

func NewHandlers(a *starter.App) *Handler {
	var err error
	Templates.Main, err = template.ParseFiles("views/main.html")
	if err != nil {
		log.Fatalf("Init template error: %v", err)
	}
	Templates.Msg, err = template.ParseFiles("views/msg.html")
	if err != nil {
		log.Fatalf("Init template error: %v", err)
	}
	Templates.Stats, err = template.ParseFiles("views/stats.html")
	if err != nil {
		log.Fatalf("Init template error: %v", err)
	}
	Templates.StatsIP, err = template.ParseFiles("views/stats-ip.html")
	if err != nil {
		log.Fatalf("Init template error: %v", err)
	}

	router := mux.NewRouter()
	h := &Handler{
		router,
		a,
	}
	router.HandleFunc("/", h.MainPage).Methods(http.MethodGet)
	router.HandleFunc("/", h.CreateLink).Methods(http.MethodPost)
	router.HandleFunc("/{link}", h.ReadLink)
	router.HandleFunc("/{link}/stats", h.LinkStats)
	router.HandleFunc("/{link}/stats/{ip}", h.LinkStatsIP)

	return h
}

func (rt *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	err := Templates.Main.ExecuteTemplate(w, "main", struct{}{})
	if err != nil {
		log.Println("MainPage", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rt *Handler) CreateLink(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("CreateLink ParseForm", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	longLink := r.PostFormValue("long_link")

	// check that the long link is not empty before create short link
	if longLink == "" {
		err := Templates.Msg.ExecuteTemplate(w, "msg", Msg{
			"link must be not empty",
			"",
			"",
			"",
			"/",
		})
		if err != nil {
			log.Println("CreateLink ExecuteTemplate", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	shortLink, err := rt.a.Ls.Create(r.Context(), longLink)
	if err != nil {
		log.Println("CreateLink Create", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Templates.Msg.ExecuteTemplate(w, "msg", Msg{
		"your link:",
		path.Join(r.Host, shortLink),
		"your link to view statistics:",
		path.Join(r.Host, shortLink, "/stats"),
		"/",
	})
	if err != nil {
		log.Println("CreateLink ExecuteTemplate", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rt *Handler) ReadLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortLink := vars["link"]

	longLink, err := rt.a.Ls.Read(r.Context(), shortLink)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("ReadLink Read", err)
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			log.Println("ReadLink Read", err)
			http.Error(w, "read link error", http.StatusInternalServerError)
		}
		return
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]
	err = rt.a.Ss.Add(r.Context(), stat.Stat{
		Link: shortLink,
		IP:   ip,
	})
	if err != nil {
		log.Println("ReadLink Add statistics error:", err)
	}

	http.Redirect(w, r, longLink, http.StatusFound)
}

func (rt *Handler) LinkStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortLink := vars["link"]

	stats, err := rt.a.Ss.ReadAll(r.Context(), shortLink)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = Templates.Msg.ExecuteTemplate(w, "msg", Msg{
				"по данному URL пока ещё никто не переходил",
				"",
				"",
				"",
				"/",
			})
			if err != nil {
				log.Println("LinkStats ReadAll", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			log.Println("LinkStats ReadAll", err)
			http.Error(w, "read statictics error", http.StatusInternalServerError)
		}
		return
	}
	err = Templates.Stats.ExecuteTemplate(w, "stats", struct{ Stats *[]stat.Stat }{Stats: stats})
	if err != nil {
		log.Println("LinkStats ExecuteTemplate", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rt *Handler) LinkStatsIP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortLink := vars["link"]
	ip := vars["ip"]

	count, err := rt.a.Ss.ReadIP(r.Context(), stat.Stat{
		Link: shortLink,
		IP:   ip,
	})
	if err != nil {
		log.Println("LinkStatsIP ReadIP", err)
		http.Error(w, "read IP statictics error", http.StatusInternalServerError)
		return
	}

	err = Templates.StatsIP.ExecuteTemplate(w, "stats-ip", struct {
		Stat  stat.Stat
		Count int64
	}{
		Stat: stat.Stat{
			Link: shortLink,
			IP:   ip,
		},
		Count: count,
	})
	if err != nil {
		log.Println("LinkStatsIP ExecuteTemplate", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
