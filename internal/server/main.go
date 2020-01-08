package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dilap54/voronov_idor/internal/app"
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
		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func New(addr string, app *app.Service) *Service {
	router := mux.NewRouter()

	router.Use(LogRequests)

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
