package controller

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/padurean/purest/internal"
	"github.com/padurean/purest/internal/auth"
	"github.com/padurean/purest/internal/database"
	"github.com/padurean/purest/internal/logging"
	"github.com/padurean/purest/internal/validator"
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
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
}

// Render ...
func (sr *SignInResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// UserRequest ...
type UserRequest struct {
	*database.User
	FirstName NullString `json:"first_name" swaggertype:"string"`
	LastName  NullString `json:"last_name" swaggertype:"string"`
	Deleted   NullTime   `json:"deleted,omitempty" swaggertype:"string"`
}

// Bind ...
func (u *UserRequest) Bind(r *http.Request) error {
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
	FirstName NullString `json:"first_name" swaggertype:"string"`
	LastName  NullString `json:"last_name" swaggertype:"string"`
	Deleted   NullTime   `json:"deleted,omitempty" swaggertype:"string"`
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
		db := r.Context().Value(internal.ContextKeyDB).(*database.DB)
		reqLogger := logging.Simple(r)
		var u *database.User
		idParam := chi.URLParam(r, "id")
		usernameOrEmail := chi.URLParam(r, "usernameOrEmail")
		if idParam == "" && usernameOrEmail == "" {
			render.Render(w, r, ErrBadRequest(
				fmt.Errorf("neither user id, nor username or email path params are specified")))
			return
		}
		if idParam != "" {
			id, err := strconv.ParseInt(idParam, 10, 64)
			if err != nil {
				render.Render(w, r, ErrBadRequest(
					fmt.Errorf("user 'id' url param '%s' is not an integer number", idParam)))
				return
			}
			u, err = (&database.User{ID: id}).GetByID(db)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					render.Render(w, r, ErrNotFound)
					return
				default:
					reqLogger.Err(err).Msgf("error getting user with id %d", id)
					render.Render(w, r, ErrInternalServer(err))
					return
				}
			}
		} else if usernameOrEmail != "" {
			var err error
			if strings.Contains(usernameOrEmail, "@") {
				u, err = (&database.User{Email: usernameOrEmail}).GetByEmail(db)
				if err != nil {
					switch err {
					case sql.ErrNoRows:
						// do nothing
					default:
						reqLogger.Err(err).Msgf("error getting user %s by email ", usernameOrEmail)
						render.Render(w, r, ErrInternalServer(err))
						return
					}
				}
			}
			if u == nil {
				u, err = (&database.User{Username: usernameOrEmail}).GetByUsername(db)
				if err != nil {
					switch err {
					case sql.ErrNoRows:
						render.Render(w, r, ErrNotFound)
						return
					default:
						reqLogger.Err(err).Msgf("error getting user %s by username", usernameOrEmail)
						render.Render(w, r, ErrInternalServer(err))
						return
					}
				}
			}
		}
		ctx := context.WithValue(r.Context(), internal.ContextKeyUser, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserCreate ...
// @id UserCreate
// @tags user
// @summary Creates a new user
// @accept application/json
// @produce application/json
// @param Authorization header string true "Bearer <token>"
// @param UserRequest body controller.UserRequest true "Request body payload"
// @success 201 {object} controller.UserResponse
// @router /users [post]
func UserCreate(w http.ResponseWriter, r *http.Request) {
	uReq := &UserRequest{}
	reqLogger := logging.Simple(r)
	if err := render.Bind(r, uReq); err != nil {
		reqLogger.Err(err).Msgf("error unmarshaling user from JSON")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	hashedPassword, err := auth.HashAndSaltPassword(uReq.Password)
	if err != nil {
		reqLogger.Err(err).Msgf("error hashing and setting password")
		render.Render(w, r, ErrUnprocessableEntity(err))
		return
	}
	uReq.Password = hashedPassword

	db := r.Context().Value(internal.ContextKeyDB).(*database.DB)
	u, err := uReq.Create(db)
	if err != nil {
		switch err.(type) {
		case *database.ErrDuplicateRow:
			render.Render(w, r, ErrUnprocessableEntity(err))
			return
		default:
			reqLogger.Err(err).Msgf("error creating user %+v", uReq.User)
			render.Render(w, r, ErrInternalServer(err))
			return
		}
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &UserResponse{User: u})
}

// UserSignIn ...
// @id UserSignIn
// @tags user
// @summary Signs-in the specified user
// @accept application/json
// @produce application/json
// @param SignInRequest body controller.SignInRequest true "Request body payload"
// @success 200 {object} controller.SignInResponse
// @router /sign-in/{usernameOrEmail} [post]
func UserSignIn(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.Simple(r)
	sReq := &SignInRequest{}
	if err := render.Bind(r, sReq); err != nil {
		reqLogger.Err(err).Msgf("error unmarshaling sign in payload from JSON")
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	u := r.Context().Value(internal.ContextKeyUser).(*database.User)
	if !auth.ComparePasswords(sReq.Password, u.Password) {
		err := fmt.Errorf("wrong password supplied for user %d", u.ID)
		reqLogger.Err(err).Msg("")
		render.Render(w, r, ErrUnauthorized(err))
		return
	}
	token, expiration, err := auth.GenerateToken(u.ID, u.Role)
	reqLogger.Debug().Msgf("generated token: %s, Error: %v", token, err)
	render.Status(r, http.StatusOK)
	render.Render(w, r, &SignInResponse{Token: token, Expiration: expiration})
}

// UserList ...
func UserList(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(internal.ContextKeyDB).(*database.DB)
	page := r.Context().Value(internal.ContextKeyPage).(int)
	pageSize := r.Context().Value(internal.ContextKeyPageSize).(int)

	u := &database.User{}
	reqLogger := logging.Simple(r)
	users, err := u.List(db, pageSize, (page-1)*pageSize)
	if err != nil {
		reqLogger.Err(err).Msgf("error listing users page")
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
	reqLogger := logging.Simple(r)
	if err := render.Bind(r, uReq); err != nil {
		reqLogger.Err(err).Msgf("error unmarshaling user from JSON")
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	u := r.Context().Value(internal.ContextKeyUser).(*database.User)
	uReq.ID = u.ID
	hashedPassword, err := auth.HashAndSaltPassword(uReq.Password)
	if err != nil {
		reqLogger.Err(err).Msgf("error hashing and setting password")
		render.Render(w, r, ErrUnprocessableEntity(err))
		return
	}
	uReq.Password = hashedPassword

	db := r.Context().Value(internal.ContextKeyDB).(*database.DB)
	u, err = uReq.Update(db)
	if err != nil {
		switch err.(type) {
		case *database.ErrDuplicateRow:
			render.Render(w, r, ErrUnprocessableEntity(err))
			return
		default:
			reqLogger.Err(err).Msgf("error updating user %+v", uReq.User)
			render.Render(w, r, ErrInternalServer(err))
			return
		}
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &UserResponse{User: u})
}

// UserGet ...
func UserGet(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(internal.ContextKeyUser).(*database.User)
	render.Status(r, http.StatusOK)
	render.Render(w, r, &UserResponse{User: u})
}

// UserDelete ...
func UserDelete(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(internal.ContextKeyUser).(*database.User)
	db := r.Context().Value(internal.ContextKeyDB).(*database.DB)
	reqLogger := logging.Simple(r)

	if err := u.Delete(db); err != nil {
		reqLogger.Err(err).Msgf("error deleting user %d", u.ID)
		render.Render(w, r, ErrUnprocessableEntity(err))
		return
	}
}
