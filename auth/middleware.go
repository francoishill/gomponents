package auth

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/francoishill/gomponents/rendering"
	"github.com/francoishill/gomponents/token"
	"github.com/francoishill/gomponents/user"
)

type Middleware interface {
	Authenticate() []func(http.Handler) http.Handler
	LoadUser() func(http.Handler) http.Handler
	GetContextUser(ctx context.Context) user.User
}

func DefaultMiddleware(userRepoFactory user.RepoFactory, rendering rendering.Service, token token.Service) *defaultMiddleware {
	type ctxKey struct{}
	return &defaultMiddleware{
		userRepoFactory,
		rendering,
		token,
		&ctxKey{},
	}
}

type defaultMiddleware struct {
	userRepoFactory user.RepoFactory
	rendering       rendering.Service
	token           token.Service

	authUserCtxKey interface{}
}

func (m *defaultMiddleware) Authenticate() []func(http.Handler) http.Handler {
	return m.token.Middlewares()
}

func (m *defaultMiddleware) LoadUser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := m.token.UserIDFromContext(r.Context())
			if err != nil {
				userMessage := "Failed to get user ID from context"
				logrus.WithError(err).Error(userMessage)
				m.rendering.RenderError(w, r, errors.Errorf(userMessage), nil, http.StatusUnauthorized)
				return
			}

			user, err := m.userRepoFactory.Repo().Get(userID)
			if err != nil {
				m.rendering.RenderError(w, r, errors.Wrapf(err, "Failed to get user"), nil, http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), m.authUserCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m *defaultMiddleware) GetContextUser(ctx context.Context) user.User {
	return ctx.Value(m.authUserCtxKey).(user.User)
}
