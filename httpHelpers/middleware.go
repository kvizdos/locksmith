package httpHelpers

import (
	"context"
	"net/http"

	"kv.codes/locksmith/database"
)

func InjectDatabaseIntoContext(next http.Handler, db database.DatabaseAccessor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Inject the custom data into the request
		r = r.WithContext(context.WithValue(r.Context(), "database", db))

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
