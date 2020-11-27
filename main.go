package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nnqq/scr-api/config"
	"github.com/nnqq/scr-api/logger"
	"github.com/nnqq/scr-api/middleware"
	"github.com/nnqq/scr-api/ratelimit"
	"github.com/nnqq/scr-proto/codegen/go/billing"
	"github.com/nnqq/scr-proto/codegen/go/category"
	"github.com/nnqq/scr-proto/codegen/go/city"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"github.com/nnqq/scr-proto/codegen/go/technology"
	"github.com/nnqq/scr-proto/codegen/go/user"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

func serveHealthz(mux *http.ServeMux) {
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write(nil)
		logger.Err(err)
	}))
}

func serveSwaggerUI(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir("./docs"))
	mux.Handle("/docs/", http.StripPrefix("/docs/", fs))
}

func serveGW(mux *runtime.ServeMux) {
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	logger.Must(parser.RegisterCompanyHandlerFromEndpoint(ctx, mux, config.Env.Service.Parser, opts))
	logger.Must(parser.RegisterPostHandlerFromEndpoint(ctx, mux, config.Env.Service.Parser, opts))
	logger.Must(city.RegisterCityHandlerFromEndpoint(ctx, mux, config.Env.Service.City, opts))
	logger.Must(category.RegisterCategoryHandlerFromEndpoint(ctx, mux, config.Env.Service.Category, opts))
	logger.Must(technology.RegisterTechnologyHandlerFromEndpoint(ctx, mux, config.Env.Service.Technology, opts))
	logger.Must(user.RegisterUserHandlerFromEndpoint(ctx, mux, config.Env.Service.User, opts))
	logger.Must(user.RegisterRoleHandlerFromEndpoint(ctx, mux, config.Env.Service.User, opts))
	logger.Must(billing.RegisterBillingHandlerFromEndpoint(ctx, mux, config.Env.Service.Billing, opts))
}

func main() {
	mux := http.NewServeMux()
	serveHealthz(mux)
	serveSwaggerUI(mux)

	gwMux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}))
	serveGW(gwMux)

	mux.Handle("/", gwMux)
	addr := strings.Join([]string{"0.0.0.0", config.Env.HTTP.Port}, ":")
	logger.Must(http.ListenAndServe(addr, middleware.CORS(ratelimit.Middleware(middleware.Auth(mux)))))
}
