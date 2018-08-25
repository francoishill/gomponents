package anonymous

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/francoishill/gomponents/auth"
	"github.com/francoishill/gomponents/rendering"
	"github.com/francoishill/gomponents/request"
)

func Router(rendering rendering.Service, auth auth.Service, requestFactory RequestFactory, responseFactory ResponseFactory) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		body := requestFactory.Register()
		if err := request.DecodeAndValidateJSON(r.Body, body); err != nil {
			rendering.RenderError(w, r, err, nil, http.StatusBadRequest)
			return
		}

		user := body.ToUser()
		token, err := auth.Register(user)
		if err != nil {
			rendering.RenderError(w, r, err, nil, http.StatusUnauthorized)
			return
		}

		render.Respond(w, r, responseFactory.LoggedIn(user, token))
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		body := requestFactory.Login()
		if err := request.DecodeAndValidateJSON(r.Body, body); err != nil {
			rendering.RenderError(w, r, err, nil, http.StatusBadRequest)
			return
		}

		user := body.ToUser()
		token, err := auth.Login(user, body.Password())
		if err != nil {
			rendering.RenderError(w, r, err, nil, http.StatusUnauthorized)
			return
		}

		render.Respond(w, r, responseFactory.LoggedIn(user, token))
	})

	r.Post("/magic-login", func(w http.ResponseWriter, r *http.Request) {
		body := requestFactory.MagicLogin()
		if err := request.DecodeAndValidateJSON(r.Body, body); err != nil {
			rendering.RenderError(w, r, err, nil, http.StatusBadRequest)
			return
		}

		user := body.ToUser()
		token, err := auth.MagicLogin(user, body.Token())
		if err != nil {
			rendering.RenderError(w, r, err, nil, http.StatusUnauthorized)
			return
		}

		render.Respond(w, r, responseFactory.LoggedIn(user, token))
	})

	return r
}
