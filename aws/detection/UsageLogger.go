package detection

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func UsageLogger(db *pgx.Conn, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		endpoint := r.URL.Path

		go func() {
			_, err := db.Exec(context.Background(),
				`INSERT INTO usage_logs (user_id, endpoint) VALUES ($1, $2)`,
				userID, endpoint,
			)
			if err != nil {

			}
		}()

		next.ServeHTTP(w, r)
	})
}
