package admin

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/pkg/errors"

	"github.com/francoishill/gomponents/auth"
	"github.com/francoishill/gomponents/encryption"
	"github.com/francoishill/gomponents/rendering"
	"github.com/francoishill/gomponents/request"
	"github.com/francoishill/gomponents/user"
)

func Router(
	authMiddleware auth.Middleware, adminMiddlware Middleware,
	requestFactory RequestFactory, responseFactory ResponseFactory,
	rendering rendering.Service,
	userRepoFactory user.RepoFactory,
	encryption encryption.Service) *chi.Mux {

	r := chi.NewRouter()

	r.Use(authMiddleware.Authenticate()...)
	r.Use(authMiddleware.LoadUser())
	r.Use(adminMiddlware.RequireAdmin())

	r.Route("/users", func(r chi.Router) {
		//list
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			users, err := userRepoFactory.Repo().List()
			if err != nil {
				rendering.RenderError(w, r, errors.Wrapf(err, "Failed to list users"), nil, http.StatusInternalServerError)
				return
			}
			responses := []UserResponse{}
			for _, u := range users {
				responses = append(responses, responseFactory.User(u))
			}
			render.Respond(w, r, responses)
		})

		//add
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			body := requestFactory.AddUser()
			if err := request.DecodeAndValidateJSON(r.Body, body); err != nil {
				rendering.RenderError(w, r, err, nil, http.StatusBadRequest)
				return
			}

			randomPassword, err := encryption.NewRandomPassword()
			if err != nil {
				rendering.RenderError(w, r, errors.Wrapf(err, "Failed to generate password"), nil, http.StatusInternalServerError)
				return
			}
			passwordHash, err := encryption.HashPassword(randomPassword)
			if err != nil {
				rendering.RenderError(w, r, errors.Wrapf(err, "Failed to hash new password"), nil, http.StatusInternalServerError)
				return
			}

			newUser := body.ToUser(passwordHash)
			if err := userRepoFactory.Repo().Add(newUser); err != nil {
				rendering.RenderError(w, r, errors.Wrapf(err, "Failed to add user"), nil, http.StatusInternalServerError)
				return
			}
			render.Respond(w, r, responseFactory.User(newUser))
		})
	})

	return r
}
