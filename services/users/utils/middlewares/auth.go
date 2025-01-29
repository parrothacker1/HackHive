package middlewares

import "net/http"

func AuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    handler.ServeHTTP(w,r)
  }
}
