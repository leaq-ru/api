package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nnqq/scr-api/config"
	"github.com/nnqq/scr-api/logger"
	"github.com/nnqq/scr-api/ratelimit"
	"github.com/nnqq/scr-proto/codegen/go/category"
	"github.com/nnqq/scr-proto/codegen/go/city"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		OrigName: false,
	}))
	opts := []grpc.DialOption{grpc.WithInsecure()}

	logger.Must(parser.RegisterCompanyHandlerFromEndpoint(ctx, mux, config.Env.Service.Parser, opts))
	logger.Must(city.RegisterCityHandlerFromEndpoint(ctx, mux, config.Env.Service.City, opts))
	logger.Must(category.RegisterCategoryHandlerFromEndpoint(ctx, mux, config.Env.Service.Category, opts))

	logger.Must(http.ListenAndServe(strings.Join([]string{
		"0.0.0.0",
		config.Env.HTTP.Port,
	}, ":"), ratelimit.Middleware.Handler(mux)))
}
