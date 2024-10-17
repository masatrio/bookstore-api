package middleware

import (
	"log"
	"net/http"
	"runtime"
)

// PanicRecoveryMiddleware recovers from panics and writes a 500 if there was one.
func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("Panic recovered: %v\n", rec)

				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, true)
				log.Printf("Stack trace:\n%s\n", buf[:stackSize])

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
