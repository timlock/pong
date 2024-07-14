package auth

import (
	"context"
	"log"
	"net/http"
	"pong/repository"

	"github.com/google/uuid"
)

type UserKey struct{}

func BasicAuth(next http.Handler, userStore *repository.UserStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, pwd, ok := req.BasicAuth()
		if !ok {
			log.Printf("User %v did not provide credentials", req.RemoteAddr)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		uuid, err := uuid.Parse(pwd)
		if err != nil {
			log.Printf("Could not parse uuid from user %v : %v", req.RemoteAddr, err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		user, ok := userStore.Get(uuid)
		if !ok {
			log.Printf("Declined user %v", req.RemoteAddr)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return

		}

		ctx := context.WithValue(context.Background(), UserKey{}, user)
		req = req.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
