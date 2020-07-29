PROTO_VER     := v1.0.0
SWAGGER_PATH  := /Users/denis/go/pkg/mod/github.com/nnqq/scr-proto@$(PROTO_VER)/codegen/swagger
COMPANY_PATH  := $(SWAGGER_PATH)/parser/company.swagger.json
CITY_PATH     := $(SWAGGER_PATH)/city/city.swagger.json
CATEGORY_PATH := $(SWAGGER_PATH)/category/category.swagger.json

all:
	go mod download;
	swagger mixin $(COMPANY_PATH) $(CITY_PATH) $(CATEGORY_PATH) -o docs/swagger.json; # https://github.com/go-swagger/go-swagger
