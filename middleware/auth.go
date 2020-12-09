package middleware

import (
	"context"
	"github.com/nnqq/scr-api/call"
	"github.com/nnqq/scr-api/logger"
	"github.com/nnqq/scr-proto/codegen/go/user"
	"net/http"
	"strings"
	"time"
)

const HeaderUserID = "Grpc-Metadata-User-Id"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.Header.Set(HeaderUserID, "")

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token != "" {
			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
			defer cancel()

			authUser, err := call.User.Auth(ctx, &user.AuthRequest{
				Token: token,
			})
			if err != nil {
				logger.Log.Error().Err(err).Send()
				w.WriteHeader(http.StatusUnauthorized)

				_, err = w.Write(nil)
				if err != nil {
					logger.Log.Error().Err(err).Send()
				}
				return
			}

			r.Header.Set(HeaderUserID, authUser.GetId())
		}

		next.ServeHTTP(w, r)
	})
}
