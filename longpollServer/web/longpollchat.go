package web

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"io"
	"longpollServer/authorization"
	"longpollServer/chat"
	"net/http"
)

type LongpollChat struct {
	chatController *chat.Controller
	authController *authorization.Controller
	*chi.Mux
	tokenAuth *jwtauth.JWTAuth
}

func NewLongpollChat(
	chatController *chat.Controller,
	authController *authorization.Controller,
	tokenAuth *jwtauth.JWTAuth,
) *LongpollChat {
	return &LongpollChat{
		chatController: chatController,
		authController: authController,
		Mux:            chi.NewMux(),
		tokenAuth:      tokenAuth,
	}
}

func (lc *LongpollChat) Start() {
	lc.Use(middleware.Logger)

	lc.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	lc.Route("/auth", func(r chi.Router) {
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			data := &LoginRequest{}
			all, err := io.ReadAll(r.Body)
			if printErr(err, w) {
				return
			}
			err = json.Unmarshal(all, &data)
			if printErr(err, w) {
				return
			}

			userId, s, err := lc.authController.Login(data.Email, data.Password)
			if printErr(err, w) {
				return
			}
			response, err := json.Marshal(
				&LoginResponse{
					UserId: userId,
					Token:  s,
				})
			if printErr(err, w) {
				return
			}
			w.Write(response)
		})
		r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			data := &RegisterRequest{}
			all, err := io.ReadAll(r.Body)
			if printErr(err, w) {
				return
			}
			err = json.Unmarshal(all, &data)
			if printErr(err, w) {
				return
			}

			userId, s, err := lc.authController.CreateUser(data.Username, data.Email, data.Password)
			if printErr(err, w) {
				return
			}
			response, err := json.Marshal(
				&RegisterResponse{
					UserId: userId,
					Token:  s,
				})
			if printErr(err, w) {
				return
			}
			w.Write(response)
		})
	})

	lc.Route("/chat", func(r chi.Router) {
		r.Use(jwtauth.Verifier(lc.tokenAuth), jwtauth.Authenticator(lc.tokenAuth))

		r.Get("/get-message", func(w http.ResponseWriter, r *http.Request) {
			message, err := lc.chatController.GetMessage(getIdFromContext(r))
			if printErr(err, w) {
				return
			}
			response, err := json.Marshal(message)
			if printErr(err, w) {
				return
			}
			w.Write(response)
		})
		r.Get("/get-messages", func(w http.ResponseWriter, r *http.Request) {
			messages, err := lc.chatController.GetMessages(getIdFromContext(r))
			if printErr(err, w) {
				return
			}
			response, err := json.Marshal(messages)
			if printErr(err, w) {
				return
			}
			w.Write(response)
		})
		r.Post("/send", func(w http.ResponseWriter, r *http.Request) {
			data := &MessageRequest{}
			all, err := io.ReadAll(r.Body)
			if printErr(err, w) {
				return
			}
			err = json.Unmarshal(all, &data)
			if printErr(err, w) {
				return
			}
			msg := chat.Message{
				From:    uuid.MustParse(getIdFromContext(r)),
				To:      uuid.MustParse(data.To),
				Content: data.Text,
			}
			err = lc.chatController.Send(msg)
			if printErr(err, w) {
				return
			}
			w.WriteHeader(http.StatusCreated)
		})
	})

	http.ListenAndServe(":3000", lc)
}

func printErr(err error, w http.ResponseWriter) bool {
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return true
	}
	return false
}

func getIdFromContext(r *http.Request) string {
	_, claims, _ := jwtauth.FromContext(r.Context())
	return claims["user_id"].(string)
}
