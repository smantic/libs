// package recoverer offers a recovery middleware that will
// output a structured log line (zap) log for panics.
package recoverer

import (
	"net/http"

	"go.uber.org/zap"
)

// Recoverer provides a recovery middleware offering a structured log line.
func Recoverer(next http.Handler, l *zap.Logger) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				l.Error("panic", zap.StackSkip("panic-trace", 1))

				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
