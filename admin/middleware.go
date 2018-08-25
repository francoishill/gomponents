package admin

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/francoishill/gomponents/auth"
	"github.com/francoishill/gomponents/rendering"
)

type Middleware interface {
	RequireAdmin() func(http.Handler) http.Handler
}

func DefaultMiddleware(rendering rendering.Service, authMiddleware auth.Middleware) *defaultMiddleware {
	return &defaultMiddleware{
		rendering,
		authMiddleware,
	}
}

type defaultMiddleware struct {
	rendering      rendering.Service
	authMiddleware auth.Middleware
}

func (m *defaultMiddleware) RequireAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := m.authMiddleware.GetContextUser(r.Context())
			if !user.IsAdmin() {
				m.rendering.RenderError(w, r, errors.Errorf("Admin permission is required for this action"), nil, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
