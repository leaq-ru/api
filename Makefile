PROTO_VER     := v1.1.1-0.20200806111930-14d63fa013a2
SWAGGER_PATH  := $(GOPATH)/pkg/mod/github.com/nnqq/scr-proto@$(PROTO_VER)/codegen/swagger
INFO_PATH     := $(SWAGGER_PATH)/swagger/info.swagger.json
COMPANY_PATH  := $(SWAGGER_PATH)/parser/company.swagger.json
CITY_PATH     := $(SWAGGER_PATH)/city/city.swagger.json
CATEGORY_PATH := $(SWAGGER_PATH)/category/category.swagger.json

# https://github.com/go-swagger/go-swagger
all:
	go mod download;
	- docker run --rm -it -e GOPATH=$(HOME)/go:/go -v $(HOME):$(HOME) -w $(shell pwd) quay.io/goswagger/swagger:v0.25.0 mixin \
		$(INFO_PATH) \
		$(COMPANY_PATH) \
		$(CITY_PATH) \
		$(CATEGORY_PATH) \
		-o docs/swagger.json;
	docker run --rm -it -e GOPATH=$(HOME)/go:/go -v $(HOME):$(HOME) -w $(shell pwd) quay.io/goswagger/swagger:v0.25.0 validate \
		docs/swagger.json; \