package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/EugeneNail/lifeline/internal/presentation/http/api"
)

// WriteJSONResponse wraps an endpoint handler, serializes its payload as JSON, and writes the result to the response writer.
func WriteJSONResponse(handler api.EndPointHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		status, payload := handler.Handle(request)

		if payload == nil {
			writer.WriteHeader(status)
			return
		}

		if err, ok := payload.(error); ok {
			http.Error(writer, err.Error(), status)
			return
		}

		var buffer bytes.Buffer
		if err := json.NewEncoder(&buffer).Encode(payload); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(status)

		if _, err := writer.Write(buffer.Bytes()); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}
}
