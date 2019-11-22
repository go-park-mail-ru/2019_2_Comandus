package middleware

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/token"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/delivery/grpc/user_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"strconv"
)

const (
	SessionName        = "user-session"
)

type Middleware struct {
	sessionStore	sessions.Store
	logger			*zap.SugaredLogger
	clientUrl		string
	token			*token.HashToken
}

func NewMiddleware(ss sessions.Store, logger *zap.SugaredLogger, token *token.HashToken, clientUrl string) Middleware{
	return Middleware{
		sessionStore: ss,
		logger:       logger,
		clientUrl:    clientUrl,
		token:        token,
	}
}

func (m *Middleware) AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.sessionStore.Get(r, SessionName)
		if err != nil {
			respond.Error(w, r, http.StatusUnauthorized, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			respond.Error(w, r, http.StatusUnauthorized, errors.New("no user_id cookie"))
			return
		}

		req := &user_grpc.UserID{
			ID:                   id.(int64),
		}
		u, err := clients.GetUserFromServer(req)
		if err != nil {
			respond.Error(w, r, http.StatusNotFound, err)
		}

		user := &model.User{
			ID:              u.ID,
			FirstName:       u.FirstName,
			SecondName:      u.SecondName,
			UserName:        u.UserType,
			Email:           u.Email,
			Password:        u.Password,
			EncryptPassword: u.EncryptPassword,
			Avatar:          u.Avatar,
			UserType:        u.UserType,
			FreelancerId:    u.FreelancerId,
			HireManagerId:   u.HireManagerId,
			CompanyId:       u.CompanyId,
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), respond.CtxKeyUser, user)))
	})
}

func (m *Middleware) RequestIDMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := strconv.Itoa(rand.Int())
		ctx := r.Context()
		ctx = context.WithValue(ctx, "rIDKey", reqID)
		w.Header().Set("Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) AccessLogMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info(r.URL.Path,
			zap.String("method:", r.Method),
			zap.String("remote_addr:", r.RemoteAddr),
			zap.String("url:", r.URL.Path),
		)
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) CORSMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "POST,PUT,DELETE,GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,csrf-token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", m.clientUrl)
		if r.Method == http.MethodOptions{
			// TODO: http.StatusOK?
			respond.Respond(w , r , http.StatusOK, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) CheckTokenMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			sess, err := m.sessionStore.Get(r, SessionName)
			if err != nil {
				err = errors.Wrapf(err, "CheckTokenMiddleware<-sessionStore.Get :")
				respond.Error(w, r, http.StatusUnauthorized, err)
				return
			}

			isEqual, err := m.token.Check(sess, r.Header.Get("csrf-token"))
			if !isEqual {
				err = errors.New("Bad token data")
				err = errors.Wrapf(err, "CheckTokenMiddleware<-Check:")
				respond.Error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}