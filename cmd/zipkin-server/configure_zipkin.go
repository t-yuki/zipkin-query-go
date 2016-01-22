package main

import (
	"io"
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/t-yuki/zipkin-go/restapi/operations"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureAPI(api *operations.ZipkinAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.ThriftConsumer = httpkit.ConsumerFunc(func(r io.Reader, target interface{}) error {
		return errors.NotImplemented("thrift consumer has not yet been implemented")
	})
	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.GetDependenciesHandler = operations.GetDependenciesHandlerFunc(func(params operations.GetDependenciesParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetDependencies has not yet been implemented")
	})
	api.GetServicesHandler = operations.GetServicesHandlerFunc(func() middleware.Responder {
		return middleware.NotImplemented("operation .GetServices has not yet been implemented")
	})
	api.GetSpansHandler = operations.GetSpansHandlerFunc(func(params operations.GetSpansParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetSpans has not yet been implemented")
	})
	api.GetTraceTraceIDHandler = operations.GetTraceTraceIDHandlerFunc(func(params operations.GetTraceTraceIDParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetTraceTraceID has not yet been implemented")
	})
	api.GetTracesHandler = operations.GetTracesHandlerFunc(func(params operations.GetTracesParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetTraces has not yet been implemented")
	})
	api.PostSpansHandler = operations.PostSpansHandlerFunc(func(params operations.PostSpansParams) middleware.Responder {
		return middleware.NotImplemented("operation .PostSpans has not yet been implemented")
	})

	api.ServerShutdown = func() {}
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
