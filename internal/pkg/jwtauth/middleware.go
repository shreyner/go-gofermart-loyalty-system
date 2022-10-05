package jwtauth

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

var AuthCookieKey = "auth"

var TokenCtxKey = &contextKey{"token"}
var ErrTokenCtxKey = &contextKey{"errToken"}

func Verifier(key []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
			tokenString := tokenFromCookie(r)

			token, err := verifyToken(tokenString)

			ctx := newContext(r.Context(), token, err)

			next.ServeHTTP(wr, r.WithContext(ctx))
		})
	}
}

func verifyToken(tokenString string) (*JwtData, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	return parseToken(tokenString)
}

func tokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(AuthCookieKey)

	if err != nil {
		return ""
	}

	return cookie.Value
}

func JwtDataFromContext(ctx context.Context) (*JwtData, error) {
	tokenData, _ := ctx.Value(TokenCtxKey).(*JwtData)
	err, _ := ctx.Value(ErrTokenCtxKey).(error)

	return tokenData, err
}

func newContext(ctx context.Context, token *JwtData, err error) context.Context {
	ctx = context.WithValue(ctx, TokenCtxKey, token)
	ctx = context.WithValue(ctx, ErrTokenCtxKey, err)

	return ctx
}

type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return fmt.Sprintf("_pkg_jwtauth %s", c.name)
}

func Authenticator(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
			tokenData, err := JwtDataFromContext(r.Context())

			if err != nil {
				http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

				return
			}

			if tokenData == nil {
				log.Error("jwt data is not define and no error")

				http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

				return
			}

			next.ServeHTTP(wr, r)
		})
	}
}
