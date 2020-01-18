package server

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/smolvitos/idoor/internal/repository"
)

var (
	tmpl *template.Template
)

func init() {
	tmpl = template.Must(template.ParseGlob("web/templates/*"))
}

type MainViewData struct {
	Users []*repository.User
}

func (s *Service) MainPage(w http.ResponseWriter, r *http.Request) {
	users, err := s.app.FindAllUsers()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	data := MainViewData{
		Users: users,
	}
	tmpl.ExecuteTemplate(w, "index.html", data)
}

type LoginViewData struct {
	Error string
}

func (s *Service) LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func (s *Service) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		tmpl.ExecuteTemplate(w, "login.html", LoginViewData{
			Error: err.Error(),
		})
		return
	}
	log.Print(r.PostForm)
	user, err := s.app.FindUserByLogin(r.PostForm.Get("login"))
	if err != nil {
		tmpl.ExecuteTemplate(w, "login.html", LoginViewData{
			Error: err.Error(),
		})
		return
	}
	if user == nil {
		tmpl.ExecuteTemplate(w, "login.html", LoginViewData{
			Error: "Пользователь не найден",
		})
		return
	}
	if user.Password != r.PostForm.Get("password") {
		tmpl.ExecuteTemplate(w, "login.html", LoginViewData{
			Error: "Неверный пароль",
		})
		return
	}
	token := s.app.GenerateAuthToken(user)
	s.app.SaveTokenForAuth(user, token)
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

type MessageViewData struct {
	SourceName string
	Text       string
}

type MessagesViewData struct {
	Users    map[uint]*repository.User
	Messages []MessageViewData
}

func (s *Service) Messages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user1 := r.Context().Value(UserKey).(*repository.User)
	user1Id := user1.ID
	user2Id, err := strconv.ParseUint(vars["userID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user2, err := s.app.FindOneUser(uint(user2Id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	messages, err := s.app.FindMessages(uint(user1Id), uint(user2Id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := MessagesViewData{
		Users: map[uint]*repository.User{
			user1.ID: user1,
			user2.ID: user2,
		},
		Messages: make([]MessageViewData, 0),
	}
	for _, m := range messages {
		data.Messages = append(data.Messages, MessageViewData{
			Text:       m.Text,
			SourceName: data.Users[m.Source].Login,
		})
	}
	log.Print(messages)
	tmpl.ExecuteTemplate(w, "messages.html", data)
}
