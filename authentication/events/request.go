package events

import (
	"context"
	"net/http"
)

type requestContextKey struct{}

// WithRequest attaches the inbound *http.Request to ctx so that Middleware
// (or anything else publishing on the bus) can later recover it via
// RequestFromContext -- e.g. to read headers, query params, cookies, or
// route values when deciding what to attach to an Envelope.
//
// authenticator, register, and sign_out_http all call this before
// publishing, so any Middleware registered via WithMiddleware can rely on
// the request being available regardless of which flow fired the event.
func WithRequest(ctx context.Context, r *http.Request) context.Context {
	if r == nil {
		return ctx
	}
	return context.WithValue(ctx, requestContextKey{}, r)
}

// RequestFromContext returns the *http.Request previously attached with
// WithRequest, if any.
func RequestFromContext(ctx context.Context) (*http.Request, bool) {
	r, ok := ctx.Value(requestContextKey{}).(*http.Request)
	return r, ok
}
