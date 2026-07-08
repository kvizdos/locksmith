package authenticator

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/kvizdos/locksmith/api_helpers"
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
)

func (a *authorizers) ServeProviderStartAPI(w http.ResponseWriter, r *http.Request) {
	ctx := a.enrichCtx(r)

	// First, detect which handler this session will
	// utilize. Since this uses a path value to determine
	// the handler, it will only match if the path value
	// is provided and matches a known handler.
	handler, err := a.getHandler(r)
	if err != nil {
		a.log.ErrorContext(ctx, "no handler supports this request", "error", err)
		api_helpers.WriteResponse(w, api_helpers.APIResponseError{
			Reason: "no handler supports this request",
		}, http.StatusInternalServerError)
		return
	}

	ctx = context.WithValue(ctx, "login_handler", handler.Name())
	ctx = context.WithValue(ctx, "login_passwordless", handler.Passwordless())
	ctx = context.WithValue(ctx, "log", a.log)

	b, ok := handler.(authenticator_domain.Beginnable)

	if !ok {
		a.log.ErrorContext(ctx, "handler does not support begin", "handler", handler.Name(), "handler_type", fmt.Sprintf("%T", handler))
		api_helpers.WriteResponse(w, api_helpers.APIResponseError{
			Reason: "no handler supports this request",
		}, http.StatusInternalServerError)
		return
	}

	if err := b.Begin(ctx, w, r); err != nil {
		if errors.Is(err, authenticator_domain.ErrMethodNotSupported) {
			a.log.DebugContext(ctx, "method not supported", "error", err, "stage", "begin")
			api_helpers.WriteResponse(w, api_helpers.APIResponseError{
				Reason: "method not supported",
			}, http.StatusMethodNotAllowed)
			return
		}
		a.log.ErrorContext(ctx, "failed to begin handler", "error", err)
		api_helpers.WriteResponse(w, api_helpers.APIResponseError{
			Reason: "something went wrong",
		}, http.StatusInternalServerError)
		return
	}
}
