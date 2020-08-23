package context

import (
	"context"
	"fmt"

	"github.com/padurean/purest/internal/auth"
	"github.com/padurean/purest/internal/database"
)

// Key ...
type Key string

// ContextKey ...
const (
	KeyDB        Key = "db"
	KeyUser      Key = "user"
	KeyJSONToken Key = "jsonToken"
	KeyPage      Key = "page"
	KeyPageSize  Key = "pageSize"
)

// Str ...
func (k Key) Str() string {
	return string(k)
}

// DB retrieves the DB handle from the given context
func DB(ctx context.Context) (*database.DB, error) {
	db, ok := ctx.Value(KeyDB).(*database.DB)
	if !ok {
		return nil, fmt.Errorf("no DB handle found in given context for key %v", KeyDB)
	}
	return db, nil
}

// User retrieves the User from the given context
func User(ctx context.Context) (*database.User, error) {
	u, ok := ctx.Value(KeyUser).(*database.User)
	if !ok {
		return nil, fmt.Errorf("no User found in given context for key %v", KeyUser)
	}
	return u, nil
}

// JSONToken retrieves the JSONToken from the given context
func JSONToken(ctx context.Context) (*auth.JSONToken, error) {
	jt, ok := ctx.Value(KeyJSONToken).(*auth.JSONToken)
	if !ok {
		return nil, fmt.Errorf("no JSONToken found in given context for key %v", KeyJSONToken)
	}
	return jt, nil
}

// Page retrieves the page number from the given context
func Page(ctx context.Context) (int, error) {
	page, ok := ctx.Value(KeyPage).(int)
	if !ok {
		return 0, fmt.Errorf("no page number found in given context for key %v", KeyPage)
	}
	return page, nil
}

// PageSize retrieves the page size from the given context
func PageSize(ctx context.Context) (int, error) {
	pageSize, ok := ctx.Value(KeyPageSize).(int)
	if !ok {
		return 0, fmt.Errorf("no page size found in given context for key %v", KeyPageSize)
	}
	return pageSize, nil
}
