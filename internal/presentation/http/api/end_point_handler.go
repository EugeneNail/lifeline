package api

import "net/http"

// EndPointHandler is the internal HTTP handler contract used by middleware and transport adapters.
type EndPointHandler interface {
	Handle(request *http.Request) (status int, payload any)
}
