package controllers

import (
	"database/sql"
	"errors"
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
	return validator.Validate(u)
}

// GetUser ...
func (u *UserRequest) GetUser() *database.User {
	uu := u.User
	uu.FirstName = sql.NullString(u.FirstName)
	uu.LastName = sql.NullString(u.LastName)
	return uu
}

// UserResponse ...
type UserResponse struct {
	*database.User
	FirstName NullString `json:"first_name"`
	LastName  NullString `json:"last_name"`
}

// Render ...
func (u *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// FromUser ...
func (u *UserResponse) FromUser(uu *database.User) {
	u.User = uu
	u.FirstName = NullString(uu.FirstName)
	u.LastName = NullString(uu.LastName)
}

// UserCreate ...
func UserCreate(w http.ResponseWriter, r *http.Request) {
	u := &UserRequest{}
	if err := render.Bind(r, u); err != nil {
		logging.Simple(r).Err(err).Msgf("error unmarshaling user from JSON")
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	db := r.Context().Value("db").(*database.DB)
	uu, err := u.GetUser().Create(db)
	if err != nil {
		logging.Simple(r).Err(err).Msgf("error inserting user %+v", u.User)
		render.Render(w, r, ErrRender(err))
		return
	}
	uur := &UserResponse{}
	uur.FromUser(uu)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, uur)
}

// UserGet ...
func UserGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("user id path param must be an integer number")))
		return
	}
	db := r.Context().Value("db").(*database.DB)
	u := database.User{ID: id}
	uu, err := u.Get(db)
	if err != nil {
		logging.Simple(r).Err(err).Msgf("error getting user with id %d", u.ID)
		render.Render(w, r, ErrRender(err))
		return
	}
	uur := &UserResponse{}
	uur.FromUser(uu)

	render.Status(r, http.StatusOK)
	render.Render(w, r, uur)
}
