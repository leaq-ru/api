package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/leaq-ru/api/config"
	"github.com/leaq-ru/api/logger"
	"github.com/leaq-ru/api/middleware"
	"github.com/leaq-ru/api/ratelimit"
	"github.com/leaq-ru/proto/codegen/go/billing"
	"github.com/leaq-ru/proto/codegen/go/exporter"
	"github.com/leaq-ru/proto/codegen/go/org"
	"github.com/leaq-ru/proto/codegen/go/parser"
	"github.com/leaq-ru/proto/codegen/go/user"
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
	logger.Must(parser.RegisterReviewHandlerFromEndpoint(ctx, mux, config.Env.Service.Parser, opts))
	logger.Must(parser.RegisterCityHandlerFromEndpoint(ctx, mux, config.Env.Service.Parser, opts))
	logger.Must(parser.RegisterCategoryHandlerFromEndpoint(ctx, mux, config.Env.Service.Parser, opts))
	logger.Must(parser.RegisterTechnologyHandlerFromEndpoint(ctx, mux, config.Env.Service.Parser, opts))
	logger.Must(parser.RegisterDnsHandlerFromEndpoint(ctx, mux, config.Env.Service.Parser, opts))
	logger.Must(user.RegisterUserHandlerFromEndpoint(ctx, mux, config.Env.Service.User, opts))
	logger.Must(user.RegisterRoleHandlerFromEndpoint(ctx, mux, config.Env.Service.User, opts))
	logger.Must(billing.RegisterBillingHandlerFromEndpoint(ctx, mux, config.Env.Service.Billing, opts))
	logger.Must(exporter.RegisterExporterHandlerFromEndpoint(ctx, mux, config.Env.Service.Exporter, opts))
	logger.Must(org.RegisterOrgHandlerFromEndpoint(ctx, mux, config.Env.Service.Org, opts))
}

func main() {
	mux := http.NewServeMux()
	serveHealthz(mux)
	serveSwaggerUI(mux)

	gwMux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}))
	serveGW(gwMux)

	mux.Handle("/", gwMux)
	addr := strings.Join([]string{"0.0.0.0", config.Env.HTTP.Port}, ":")
	logger.Must(http.ListenAndServe(addr, middleware.CORS(middleware.Auth(ratelimit.Middleware(mux)))))
}
