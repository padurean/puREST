package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/padurean/purest/database"
	"github.com/padurean/purest/logging"
	"github.com/padurean/purest/validator"
)

// ContextKey ...
const (
	ContextKeyDB       ContextKey = "db"
	ContextKeyUser     ContextKey = "user"
	ContextKeyPage     ContextKey = "page"
	ContextKeyPageSize ContextKey = "pageSize"
	// ...
)

// SignInRequest ...
type SignInRequest struct {
	Password string `json:"password" validate:"required"`
}

// Bind ...
func (sr *SignInRequest) Bind(r *http.Request) error {
	if err := validator.Validate(sr); err != nil {
		return err
	}
	return nil
}

// SignInResponse ...
type SignInResponse struct {
	Token string `json:"token"`
}

// Render ...
func (sr *SignInResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// UserRequest ...
type UserRequest struct {
	*database.User
	FirstName NullString `json:"first_name"`
	LastName  NullString `json:"last_name"`
	Deleted   NullTime   `json:"deleted,omitempty"`
}

// Bind ...
func (u *UserRequest) Bind(r *http.Request) error {
	// TODO OGG: add validation for username (use the validator with a custom tag/validation maybe?)
	if err := validator.Validate(u); err != nil {
		return err
	}
	u.User.FirstName = sql.NullString(u.FirstName)
	u.User.LastName = sql.NullString(u.LastName)
	u.User.Deleted = sql.NullTime(u.Deleted)
	return nil
}

// UserResponse ...
type UserResponse struct {
	*database.User
	FirstName NullString `json:"first_name"`
	LastName  NullString `json:"last_name"`
	Deleted   NullTime   `json:"deleted,omitempty"`
}

// Render ...
func (u *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	u.FirstName = NullString(u.User.FirstName)
	u.LastName = NullString(u.User.LastName)
	u.Deleted = NullTime(u.User.Deleted)
	return nil
}

// UserCtx ...
func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := r.Context().Value(ContextKeyDB).(*database.DB)

		if idParam := chi.URLParam(r, "id"); idParam != "" {
			// get user by id
			id, err := strconv.ParseInt(idParam, 10, 64)
			if err != nil {
				render.Render(w, r, ErrBadRequest(fmt.Errorf("user 'id' url param '%s' is not an integer number", idParam)))
				return
			}
			u, err := (&database.User{ID: id}).GetByID(db)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
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
		} else if usernameOrEmail := chi.URLParam(r, "usernameOrEmail"); usernameOrEmail != "" {
			u := &database.User{Username: usernameOrEmail, Email: usernameOrEmail}
			var uu *database.User
			var err error
			// TODO OGG: if usernameOrEmail contains @ fetch by email and if not found, by username
			//					 else fetch only by username
			uu, err = u.GetByUsername(db)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					uu, err = u.GetByEmail(db)
					if err != nil {
						switch err {
						case sql.ErrNoRows:
							render.Render(w, r, ErrNotFound)
							return
						default:
							logging.Simple(r).Err(err).Msgf("error getting user with email %d", u.Email)
							render.Render(w, r, ErrInternalServer(err))
							return
						}
					}
				default:
					logging.Simple(r).Err(err).Msgf("error getting user with username %d", u.Username)
					render.Render(w, r, ErrInternalServer(err))
					return
				}
			}
			ctx := context.WithValue(r.Context(), ContextKeyUser, uu)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			render.Render(w, r, ErrBadRequest(fmt.Errorf("username or email url param not specified")))
			return
		}
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

	if err := uReq.HashSaltAndSetPassword(); err != nil {
		logging.Simple(r).Err(err).Msgf("error hashing and setting password")
		render.Render(w, r, ErrUnprocessableEntity(err))
		return
	}

	db := r.Context().Value(ContextKeyDB).(*database.DB)
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

// UserSignIn ...
func UserSignIn(w http.ResponseWriter, r *http.Request) {
	sReq := &SignInRequest{}
	if err := render.Bind(r, sReq); err != nil {
		logging.Simple(r).Err(err).Msgf("error unmarshaling sign in payload from JSON")
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	u := r.Context().Value(ContextKeyUser).(*database.User)
	if !u.ComparePasswords(sReq.Password) {
		logging.Simple(r).Error().Msgf("wrong password supplied for user %d", u.ID)
		render.Render(w, r, ErrWrongPassword)
		return
	}
	// TODO OGG: generate token using go-chi/jwtauth at first, then try to switch to paseto
	render.Status(r, http.StatusOK)
	render.Render(w, r, &SignInResponse{Token: ""})
}

// UserList ...
func UserList(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(ContextKeyDB).(*database.DB)
	page := r.Context().Value(ContextKeyPage).(int)
	pageSize := r.Context().Value(ContextKeyPageSize).(int)

	u := &database.User{}
	users, err := u.List(db, pageSize, (page-1)*pageSize)
	if err != nil {
		logging.Simple(r).Err(err).Msgf("error listing users page")
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	usersResponseList := []render.Renderer{}
	for _, u := range users {
		usersResponseList = append(usersResponseList, &UserResponse{User: u})
	}
	if err := render.RenderList(w, r, usersResponseList); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Status(r, http.StatusOK)
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
	if err := uReq.HashSaltAndSetPassword(); err != nil {
		logging.Simple(r).Err(err).Msgf("error hashing and setting password")
		render.Render(w, r, ErrUnprocessableEntity(err))
		return
	}

	db := r.Context().Value(ContextKeyDB).(*database.DB)
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

// UserDelete ...
func UserDelete(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(ContextKeyUser).(*database.User)
	db := r.Context().Value(ContextKeyDB).(*database.DB)

	if err := u.Delete(db); err != nil {
		logging.Simple(r).Err(err).Msgf("error deleting user %s", u.ID)
		render.Render(w, r, ErrUnprocessableEntity(err))
		return
	}
}
