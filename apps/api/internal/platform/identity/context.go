package identity

import (
	"context"
	"net/http"
	"strings"
)

const (
	OwnerTypeHeader = "X-Owner-Type"
	OwnerIDHeader   = "X-Owner-ID"
	ActorIDHeader   = "X-Actor-ID"
)

type Principal struct {
	OwnerType string
	OwnerID   string
	ActorID   string
}

type contextKey struct{}

func DevelopmentHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal := Principal{
			OwnerType: strings.ToUpper(strings.TrimSpace(r.Header.Get(OwnerTypeHeader))),
			OwnerID:   strings.TrimSpace(r.Header.Get(OwnerIDHeader)),
			ActorID:   strings.TrimSpace(r.Header.Get(ActorIDHeader)),
		}
		if principal.ActorID == "" {
			principal.ActorID = principal.OwnerID
		}
		ctx := context.WithValue(r.Context(), contextKey{}, principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func FromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(contextKey{}).(Principal)
	return principal, ok && principal.OwnerType != "" && principal.OwnerID != "" && principal.ActorID != ""
}
