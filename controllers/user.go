package controllers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/padurean/purest/database"
	"github.com/padurean/purest/logging"
	"github.com/padurean/purest/validator"
)

// UserRequest ...
type UserRequest struct {
	*database.User
	FirstName NullString `json:"first_name"`
	LastName  NullString `json:"last_name"`
}

// Bind ...
func (u *UserRequest) Bind(r *http.Request) error {
	if err := validator.Validate(u); err != nil {
		return err
	}
	u.User.FirstName = sql.NullString(u.FirstName)
	u.User.LastName = sql.NullString(u.LastName)
	return nil
}

// UserResponse ...
type UserResponse struct {
	*database.User
	FirstName NullString `json:"first_name"`
	LastName  NullString `json:"last_name"`
}

// Render ...
func (u *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	u.FirstName = NullString(u.User.FirstName)
	u.LastName = NullString(u.User.LastName)
	return nil
}

// UserCtx ...
func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := r.Context().Value("db").(*database.DB)

		if idParam := chi.URLParam(r, "id"); idParam != "" {
			// get user by id
			id, err := strconv.ParseInt(idParam, 10, 64)
			if err != nil {
				render.Render(w, r, ErrBadRequest(errors.New("user id url param must be an integer number")))
				return
			}
			u, err := (&database.User{ID: id}).GetByID(db)
			if err != nil {
				switch {
				case err == sql.ErrNoRows:
					render.Render(w, r, ErrNotFound)
					return
				default:
					logging.Simple(r).Err(err).Msgf("error getting user with id %d", id)
					render.Render(w, r, ErrInternalServer(err))
					return
				}
			}
			ctx := context.WithValue(r.Context(), ContextKeyUser, u)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			render.Render(w, r, ErrBadRequest(fmt.Errorf("user id url param not specified")))
			return
		}
		// TODO OGG: getting user by username or email go here (different URL param names)
		// else if {
		//
		// }
	})
}

// UserCreate ...
func UserCreate(w http.ResponseWriter, r *http.Request) {
	uReq := &UserRequest{}
	if err := render.Bind(r, uReq); err != nil {
		logging.Simple(r).Err(err).Msgf("error unmarshaling user from JSON")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	db := r.Context().Value("db").(*database.DB)
	u, err := uReq.Create(db)
	if err != nil {
		switch err.(type) {
		case *database.ErrDuplicateRow:
			render.Render(w, r, ErrUnprocessableEntity(err))
			return
		default:
			logging.Simple(r).Err(err).Msgf("error creating user %+v", uReq.User)
			render.Render(w, r, ErrInternalServer(err))
			return
		}
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &UserResponse{User: u})
}

// UserUpdate ...
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	uReq := &UserRequest{}
	if err := render.Bind(r, uReq); err != nil {
		logging.Simple(r).Err(err).Msgf("error unmarshaling user from JSON")
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	u := r.Context().Value(ContextKeyUser).(*database.User)
	uReq.ID = u.ID

	db := r.Context().Value("db").(*database.DB)
	u, err := uReq.Update(db)
	if err != nil {
		switch err.(type) {
		case *database.ErrDuplicateRow:
			render.Render(w, r, ErrUnprocessableEntity(err))
			return
		default:
			logging.Simple(r).Err(err).Msgf("error updating user %+v", uReq.User)
			render.Render(w, r, ErrInternalServer(err))
			return
		}
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &UserResponse{User: u})
}

// UserGet ...
func UserGet(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(ContextKeyUser).(*database.User)
	render.Status(r, http.StatusOK)
	render.Render(w, r, &UserResponse{User: u})
}
