PROTO_VER       := v1.23.0
SWAGGER_PATH    := $(GOPATH)/pkg/mod/github.com/leaq-ru/proto@$(PROTO_VER)/codegen/swagger
INFO_PATH       := $(SWAGGER_PATH)/swagger/info.swagger.json
COMPANY_PATH    := $(SWAGGER_PATH)/parser/company.swagger.json
POST_PATH       := $(SWAGGER_PATH)/parser/post.swagger.json
REVIEW_PATH     := $(SWAGGER_PATH)/parser/review.swagger.json
CITY_PATH       := $(SWAGGER_PATH)/parser/city.swagger.json
CATEGORY_PATH   := $(SWAGGER_PATH)/parser/category.swagger.json
TECHNOLOGY_PATH := $(SWAGGER_PATH)/parser/technology.swagger.json
DNS_PATH        := $(SWAGGER_PATH)/parser/dns.swagger.json
USER_PATH       := $(SWAGGER_PATH)/user/user.swagger.json
ROLE_PATH       := $(SWAGGER_PATH)/user/role.swagger.json
BILLING_PATH    := $(SWAGGER_PATH)/billing/billing.swagger.json
EXPORTER_PATH   := $(SWAGGER_PATH)/exporter/exporter.swagger.json
ORG_PATH        := $(SWAGGER_PATH)/org/org.swagger.json

# https://github.com/go-swagger/go-swagger
all:
	go mod download;
	- docker run --rm -it -e GOPATH=$(HOME)/go:/go -v $(HOME):$(HOME) -w $(shell pwd) quay.io/goswagger/swagger:v0.25.0 mixin \
		$(INFO_PATH) \
		$(COMPANY_PATH) \
		$(POST_PATH) \
		$(REVIEW_PATH) \
		$(CITY_PATH) \
		$(CATEGORY_PATH) \
		$(TECHNOLOGY_PATH) \
		$(DNS_PATH) \
		$(USER_PATH) \
		$(ROLE_PATH) \
		$(BILLING_PATH) \
		$(EXPORTER_PATH) \
		$(ORG_PATH) \
		-o docs/swagger.json;
	docker run --rm -it -e GOPATH=$(HOME)/go:/go -v $(HOME):$(HOME) -w $(shell pwd) quay.io/goswagger/swagger:v0.25.0 flatten \
		docs/swagger.json \
		--with-flatten=remove-unused \
		-o docs/swagger.json; \
	docker run --rm -it -e GOPATH=$(HOME)/go:/go -v $(HOME):$(HOME) -w $(shell pwd) quay.io/goswagger/swagger:v0.25.0 validate \
		docs/swagger.json; \
