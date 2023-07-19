package web

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
)

// Respond converts a Go value to JSON and sends it to the client.
func Respond(ctx context.Context, w http.ResponseWriter, r *http.Request, data any, statusCode int) {
	ctx, span := AddSpan(ctx, "foundation.web.response", attribute.Int("status", statusCode))
	defer span.End()

	SetStatusCode(ctx, statusCode)

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		render.Respond(w, r, data)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Println("error:", err)
	}
	_, _ = w.Write(b)
}
