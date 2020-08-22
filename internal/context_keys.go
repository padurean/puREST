package internal

// ContextKey ...
type ContextKey string

// ContextKey ...
const (
	ContextKeyDB        ContextKey = "db"
	ContextKeyUser      ContextKey = "user"
	ContextKeyJSONToken ContextKey = "jsonToken"
	ContextKeyPage      ContextKey = "page"
	ContextKeyPageSize  ContextKey = "pageSize"
)

// Str ...
func (k ContextKey) Str() string {
	return string(k)
}
