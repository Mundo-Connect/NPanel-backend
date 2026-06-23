package server

import (
	nethttp "net/http"

	"github.com/npanel-dev/NPanel-backend/internal/conf"
	pkgmiddleware "github.com/npanel-dev/NPanel-backend/internal/pkg/middleware"
)

func newCORSAwareFallbackHandler(corsConfig *conf.Server_CORS, statusCode int) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		pkgmiddleware.ApplyCORSHeaders(w.Header(), r, corsConfig)
		if pkgmiddleware.IsCORSPreflightRequest(r) {
			w.WriteHeader(nethttp.StatusNoContent)
			return
		}

		nethttp.Error(w, nethttp.StatusText(statusCode), statusCode)
	})
}
