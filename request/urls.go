package request

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/francoishill/gomponents/rendering"
)

func RequiredURLParam(w http.ResponseWriter, r *http.Request, key string, rendering rendering.Service) (string, bool) {
	val := strings.TrimSpace(chi.URLParam(r, key))
	if val == "" {
		tmpErr := errors.Errorf("Param %s is missing", key)
		rendering.RenderError(w, r, tmpErr, nil, http.StatusBadRequest)
		return "", false
	}
	return val, true
}

func RequiredQueryParam(w http.ResponseWriter, r *http.Request, key string, rendering rendering.Service) (string, bool) {
	val := strings.TrimSpace(r.URL.Query().Get(key))
	if val == "" {
		tmpErr := errors.Errorf("Query %s is missing", key)
		rendering.RenderError(w, r, tmpErr, nil, http.StatusBadRequest)
		return "", false
	}
	return val, true
}
