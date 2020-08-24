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
	"github.com/padurean/purest/internal/auth"
	icontext "github.com/padurean/purest/internal/context"
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
	Warning    string    `json:"warning,omitempty"`
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

// UserUpdatePasswordRequest ...
type UserUpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,password"`
}

// Bind ...
func (u *UserUpdatePasswordRequest) Bind(r *http.Request) error {
	if err := validator.Validate(u); err != nil {
		return err
	}
	return nil
}

// UserUpdateEmailRequest ...
type UserUpdateEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Bind ...
func (u *UserUpdateEmailRequest) Bind(r *http.Request) error {
	if err := validator.Validate(u); err != nil {
		return err
	}
	return nil
}

// SignedInUserCtx ...
func SignedInUserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db, err := icontext.DB(r.Context())
		if err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}
		reqLogger := logging.Simple(r)
		var u *database.User
		jsonToken, err := icontext.JSONToken(r.Context())
		if err != nil {
			render.Render(w, r, ErrUnauthorized(err))
			return
		}
		u, err = (&database.User{ID: jsonToken.UserID}).GetByID(db)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				render.Render(w, r, ErrNotFound)
				return
			default:
				reqLogger.Err(err).Msgf("error getting signed-in user with id %d", jsonToken.UserID)
				render.Render(w, r, ErrInternalServer(err))
				return
			}
		}
		ctx := context.WithValue(r.Context(), icontext.KeySignedInUser, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserCtx ...
func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db, err := icontext.DB(r.Context())
		if err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}
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
		ctx := context.WithValue(r.Context(), icontext.KeyUser, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserCreate ...
// @id UserCreate
// @tags users
// @summary Creates a new user
// @accept application/json
// @produce application/json
// @param Authorization header string true "Bearer <token>"
// @param payload body controller.UserRequest true "Request body payload"
// @success 201 {object} controller.UserResponse
// @failure 401 {object} controller.ErrResponse
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

	db, err := icontext.DB(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
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
// @tags users
// @summary Signs-in the specified user
// @accept application/json
// @produce application/json
// @param usernameOrEmail path string true "Username or email"
// @param payload body controller.SignInRequest true "Request body payload"
// @success 200 {object} controller.SignInResponse
// @failure 401 {object} controller.ErrResponse
// @failure 404 {object} controller.ErrResponse
// @router /users/sign-in/{usernameOrEmail} [post]
func UserSignIn(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.Simple(r)
	sReq := &SignInRequest{}
	if err := render.Bind(r, sReq); err != nil {
		reqLogger.Err(err).Msgf("error unmarshaling sign in payload from JSON")
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	u, err := icontext.User(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	if !auth.ComparePasswords(sReq.Password, u.Password) {
		err := fmt.Errorf("wrong password supplied for user %d", u.ID)
		reqLogger.Err(err).Msg("")
		render.Render(w, r, ErrUnauthorized(err))
		return
	}
	token, expiration, err := auth.GenerateToken(u.ID, u.Role)
	reqLogger.Debug().Msgf("generated token: %s, Error: %v", token, err)
	warning := ""
	if u.Username == auth.DefaultAdminUser && sReq.Password == auth.DefaultAdminPassword {
		warning = fmt.Sprintf(
			"%s user is using the default password, to improve security please change it ASAP",
			u.Username)
	}
	sReq.Password = ""
	render.Status(r, http.StatusOK)
	render.Render(w, r, &SignInResponse{
		Token:      token,
		Expiration: expiration,
		Warning:    warning})
}

// UserList ...
// @id UserList
// @tags users
// @summary Lists users
// @accept application/json
// @produce application/json
// @param Authorization header string true "Bearer <token>"
// @param page query int false "Page number"
// @param pageSize query int false "Page size (default 20)"
// @success 200 {array} controller.UserResponse
// @failure 401 {object} controller.ErrResponse
// @router /users [get]
func UserList(w http.ResponseWriter, r *http.Request) {
	db, err := icontext.DB(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	page, err := icontext.Page(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	pageSize, err := icontext.PageSize(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

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
// @id UserUpdate
// @tags users
// @summary Updates an existing user
// @accept application/json
// @produce application/json
// @param Authorization header string true "Bearer <token>"
// @param id path int true "User id"
// @param payload body controller.UserRequest true "Request body payload"
// @success 200 {object} controller.UserResponse
// @failure 401 {object} controller.ErrResponse
// @failure 404 {object} controller.ErrResponse
// @router /users/{id} [put]
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	uReq := &UserRequest{}
	reqLogger := logging.Simple(r)
	if err := render.Bind(r, uReq); err != nil {
		reqLogger.Err(err).Msgf("error unmarshaling user from JSON")
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	u, err := icontext.User(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	uReq.ID = u.ID
	hashedPassword, err := auth.HashAndSaltPassword(uReq.Password)
	if err != nil {
		reqLogger.Err(err).Msgf("error hashing and setting password")
		render.Render(w, r, ErrUnprocessableEntity(err))
		return
	}
	uReq.Password = hashedPassword

	db, err := icontext.DB(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
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

// UserUpdatePassword ...
// @id UserUpdatePassword
// @tags users
// @summary Updates the password for the currently signed-in user
// @accept application/json
// @produce application/json
// @param Authorization header string true "Bearer <token>"
// @param payload body controller.UserUpdatePasswordRequest true "Request body payload"
// @success 200 {object} controller.UserResponse
// @failure 401 {object} controller.ErrResponse
// @failure 404 {object} controller.ErrResponse
// @router /users/password [put]
func UserUpdatePassword(w http.ResponseWriter, r *http.Request) {
	uReq := &UserUpdatePasswordRequest{}
	reqLogger := logging.Simple(r)
	if err := render.Bind(r, uReq); err != nil {
		reqLogger.Err(err).Msgf("error unmarshaling password update from JSON")
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	u, err := icontext.SignedInUser(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	if !auth.ComparePasswords(uReq.OldPassword, u.Password) {
		render.Render(w, r, ErrUnauthorized(fmt.Errorf("old password is incorrect")))
		return
	}
	hashedPassword, err := auth.HashAndSaltPassword(uReq.NewPassword)
	if err != nil {
		reqLogger.Err(err).Msgf("error hashing and setting password")
		render.Render(w, r, ErrUnprocessableEntity(err))
		return
	}
	uReq.OldPassword = ""
	uReq.NewPassword = ""
	u.Password = hashedPassword

	db, err := icontext.DB(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	u, err = u.Update(db)
	if err != nil {
		reqLogger.Err(err).Msgf("error updating password for user %d", u.ID)
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &UserResponse{User: u})
}

// UserUpdateEmail ...
// @id UserUpdateEmail
// @tags users
// @summary Updates the email for the currently signed-in user
// @accept application/json
// @produce application/json
// @param Authorization header string true "Bearer <token>"
// @param payload body controller.UserUpdateEmailRequest true "Request body payload"
// @success 200 {object} controller.UserResponse
// @failure 401 {object} controller.ErrResponse
// @failure 404 {object} controller.ErrResponse
// @router /users/email [put]
func UserUpdateEmail(w http.ResponseWriter, r *http.Request) {
	uReq := &UserUpdateEmailRequest{}
	reqLogger := logging.Simple(r)
	if err := render.Bind(r, uReq); err != nil {
		reqLogger.Err(err).Msgf("error unmarshaling email update from JSON")
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	u, err := icontext.SignedInUser(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	u.Email = uReq.Email

	db, err := icontext.DB(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	u, err = u.Update(db)
	if err != nil {
		switch err.(type) {
		case *database.ErrDuplicateRow:
			render.Render(w, r, ErrUnprocessableEntity(err))
			return
		default:
			reqLogger.Err(err).Msgf("error updating email for user %d", u.ID)
			render.Render(w, r, ErrInternalServer(err))
			return
		}
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &UserResponse{User: u})
}

// UserGet ...
// @id UserGet
// @tags users
// @summary Gets an existing user
// @accept application/json
// @produce application/json
// @param Authorization header string true "Bearer <token>"
// @param id path int true "User id"
// @success 200 {object} controller.UserResponse
// @failure 401 {object} controller.ErrResponse
// @failure 404 {object} controller.ErrResponse
// @router /users/{id} [get]
func UserGet(w http.ResponseWriter, r *http.Request) {
	u, err := icontext.User(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, &UserResponse{User: u})
}

// UserGetMe ...
// @id UserGetMe
// @tags users
// @summary Gets the currently signed-in user
// @accept application/json
// @produce application/json
// @param Authorization header string true "Bearer <token>"
// @success 200 {object} controller.UserResponse
// @failure 401 {object} controller.ErrResponse
// @failure 404 {object} controller.ErrResponse
// @router /users/me [get]
func UserGetMe(w http.ResponseWriter, r *http.Request) {
	u, err := icontext.SignedInUser(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, &UserResponse{User: u})
}

// UserDelete ...
// @id UserDelete
// @tags users
// @summary Deletes an existing user
// @accept application/json
// @produce application/json
// @param Authorization header string true "Bearer <token>"
// @param id path int true "User id"
// @success 204
// @failure 401 {object} controller.ErrResponse
// @failure 404 {object} controller.ErrResponse
// @router /users/{id} [delete]
func UserDelete(w http.ResponseWriter, r *http.Request) {
	u, err := icontext.User(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	db, err := icontext.DB(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	reqLogger := logging.Simple(r)

	if err := u.Delete(db); err != nil {
		reqLogger.Err(err).Msgf("error deleting user %d", u.ID)
		render.Render(w, r, ErrUnprocessableEntity(err))
		return
	}
	render.Status(r, http.StatusNoContent)
}
