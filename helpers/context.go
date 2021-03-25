package helpers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	CtxSession = "AppSession"
)

func AddContext(sess *Session) mux.MiddlewareFunc {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := context.WithValue(r.Context(), CtxSession, sess)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSession(req *http.Request) (*Session, error) {
	rmc := req.Context().Value(CtxSession)
	if rmc == nil {
		return nil, nil //fmt.Errorf("unable to get session from context")
	}
	sess, ok := rmc.(*Session)
	if !ok {
		return nil, fmt.Errorf("got unexpected type for %s, unable to typeassert", CtxSession)
	}

	return sess, nil
}
