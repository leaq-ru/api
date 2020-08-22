FROM golang:1.15-alpine AS build
ARG GITHUB_TOKEN=$GITHUB_TOKEN
WORKDIR /app
COPY .. /app
ENV GOPRIVATE="github.com/nnqq/*"
RUN apk add --no-cache git
RUN git config --global url."https://nnqq:$GITHUB_TOKEN@github.com/".insteadOf "https://github.com/"
RUN go build -o servicebin

FROM alpine:latest
COPY --from=build /app/servicebin /app/
