package adsvc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func SetupRouter(router *httprouter.Router) *httprouter.Router {
	endpoint := Endpoint{}
	svc := Service{}
	router.GET(wrap("/", endpoint.Index(svc)))
	router.GET(wrap("/advertisements", endpoint.All(svc)))
	router.GET("/success", endpoint.Success())
	// router.GET("/advertisements/:id", endpoint.One(svc))
	router.POST("/advertisements", endpoint.CreateForm(svc))
	router.DELETE("/advertisements/:id", endpoint.Delete(svc))
	return router
}

func wrap(p string, h func(http.ResponseWriter, *http.Request)) (string, httprouter.Handle) {
	return p, wrapHandler(alice.New(loggerMiddleware, timeoutMiddleware).ThenFunc(h))
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}
}
