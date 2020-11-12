package limiter

import "net/http"

//Limiter interface enable switching between redis rate limiter or other kinds of rate limiter
type Limiter interface {
	Limit(http.Handler) http.Handler
}
