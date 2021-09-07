ARG GOLANG_VERSION=1.14
ARG TINI_VERSION=v0.19.0
FROM golang:${GOLANG_VERSION}-alpine as build

WORKDIR $GOPATH/src/github.com/rugwirobaker/hermes
COPY go.mod go.sum  ./
RUN GO111MODULE=on GOPROXY="https://proxy.golang.org" go mod download
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -o /bin/hermes ./cmd/hermes

FROM scratch
WORKDIR /
EXPOSE 8080

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/hermes /bin/hermes

ENTRYPOINT ["/bin/hermes"]