package apiserver

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"strconv"
)

func (s *Server) AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			general.Error(w, r, http.StatusUnauthorized, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			general.Error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.usecase.Find(id.(int64))
		if err != nil {
			general.Error(w, r, http.StatusNotFound, err)
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, &u)))
	})
}

func (s *Server) RequestIDMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := strconv.Itoa(rand.Int())
		ctx := r.Context()
		ctx = context.WithValue(ctx, "rIDKey", reqID)
		w.Header().Set("Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) AccessLogMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info(r.URL.Path,
			zap.String("method:", r.Method),
			zap.String("remote_addr:", r.RemoteAddr),
			zap.String("url:", r.URL.Path),
		)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) CORSMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "POST,PUT,DELETE,GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,X-Lol")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", s.clientUrl)
		if r.Method == http.MethodOptions{
			// TODO: http.StatusOK?
			general.Error(w , r , http.StatusOK, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) CheckTokenMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/token" {
			sess, err := s.sessionStore.Get(r, sessionName)
			if err != nil {
				err = errors.Wrapf(err, "CheckTokenMiddleware<-sessionStore.Get :")
				general.Error(w, r, http.StatusUnauthorized, err)
				return
			}

			isEqual, err := s.token.Check(sess, r.Header.Get("csrf-token"))
			if !isEqual {
				err = errors.Wrapf(err, "CheckTokenMiddleware<-Check:")
				general.Error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}