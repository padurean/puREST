package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
	"github.com/padurean/purest/internal"
	"github.com/padurean/purest/internal/auth"
	"github.com/padurean/purest/internal/controller"
)

func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			render.Render(w, r, controller.ErrUnauthorized(errors.New("missing Authorization header")))
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		jsonToken, err := auth.VerifyToken(token)
		if err != nil {
			render.Render(w, r, controller.ErrUnauthorized(err))
			return
		}
		ctx := context.WithValue(r.Context(), internal.ContextKeyJSONToken, jsonToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

const pageSizeDefault = 20

func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := 1
		var err error
		if pageParam := r.URL.Query().Get(internal.ContextKeyPage.Str()); pageParam != "" {
			page, err = strconv.Atoi(pageParam)
			if err != nil {
				render.Render(w, r, controller.ErrBadRequest(
					fmt.Errorf("'page' url param '%s' is not an integer number", pageParam)))
				return
			}
		}
		pageSize := pageSizeDefault
		if pageSizeParam := r.URL.Query().Get(internal.ContextKeyPageSize.Str()); pageSizeParam != "" {
			pageSize, err = strconv.Atoi(pageSizeParam)
			if err != nil {
				render.Render(w, r, controller.ErrBadRequest(
					fmt.Errorf("'pageSize' url param '%s' is not an integer number", pageSizeParam)))
				return
			}
		}
		ctx := context.WithValue(r.Context(), internal.ContextKeyPage, page)
		ctx = context.WithValue(ctx, internal.ContextKeyPageSize, pageSize)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
