package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chr4/pwgen"
	"github.com/gorilla/mux"

	"github.com/smolvitos/idoor/internal/app"
)

type ContextKey string

const (
	UserKey ContextKey = "user"
)

type Service struct {
	app    *app.Service
	router *mux.Router
	srv    *http.Server
}

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request %s - %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func RandomCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    "remixlang",
			Value:   "0",
			Expires: time.Now().Add(724 * time.Hour),
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "remixstid",
			Value:   fmt.Sprintf("%s_%s", pwgen.Num(10), pwgen.Alpha(20)),
			Expires: time.Now().Add(724 * time.Hour),
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "remixflash",
			Value:   "0.0.0",
			Expires: time.Now().Add(724 * time.Hour),
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "remixscreen_depth",
			Value:   "24",
			Expires: time.Now().Add(724 * time.Hour),
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "remixdt",
			Value:   "0",
			Expires: time.Now().Add(724 * time.Hour),
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "tmr_reqNum",
			Value:   pwgen.Num(2),
			Expires: time.Now().Add(724 * time.Hour),
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "tmr_lvid",
			Value:   pwgen.Alpha(20),
			Expires: time.Now().Add(724 * time.Hour),
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "tmr_lvidTS",
			Value:   fmt.Sprintf("%d", time.Now().Unix()),
			Expires: time.Now().Add(724 * time.Hour),
		})
		if c, _ := r.Cookie("remixsid"); c == nil {
			http.SetCookie(w, &http.Cookie{
				Name:    "remixsid",
				Value:   pwgen.Alpha(36),
				Expires: time.Now().Add(724 * time.Hour),
			})
			http.SetCookie(w, &http.Cookie{
				Name:    "remixusid",
				Value:   "NGE1YTNiYjVlOWNiODIxNTJjMGQzZDR",
				Expires: time.Now().Add(724 * time.Hour),
			})
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Service) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, _ := r.Cookie("token")
		if tokenCookie == nil {
			log.Printf("Auth failed: cookie not found")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user, err := s.app.FindUserByToken(tokenCookie.Value)
		if err != nil {
			log.Printf("Auth failed: %s", err.Error())
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		log.Printf("User: %v", user)
		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func New(addr string, app *app.Service) *Service {
	router := mux.NewRouter()

	router.Use(LogRequests)
	router.Use(RandomCookies)

	fileServer := http.FileServer(http.Dir("web/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	svc := &Service{
		srv: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		app: app,
	}

	router.Handle("/", svc.Auth(http.HandlerFunc(svc.MainPage))).Methods("GET")
	router.Handle("/messages/{userID}", svc.Auth(http.HandlerFunc(svc.Messages))).Methods("GET")
	router.HandleFunc("/login", svc.LoginPage).Methods("GET")
	router.HandleFunc("/login", svc.Login).Methods("POST")

	return svc
}

func (s *Service) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Service) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
