ARG GOLANG_VERSION=1.19.0
ARG TINI_VERSION=v0.19.0

FROM flyio/litefs:sha-a298335 as litefs

FROM golang:${GOLANG_VERSION}-alpine as build
WORKDIR $GOPATH/src/github.com/rugwirobaker/hermes
RUN apk add build-base
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-download GOPROXY="https://proxy.golang.org" go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=1 go build -o /bin/hermes ./cmd/hermes

FROM alpine
WORKDIR /
EXPOSE 8080

RUN apk add bash curl fuse sqlite

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/hermes /bin/hermes
COPY --from=litefs /usr/local/bin/litefs /usr/local/bin/litefs

# Copy our LiteFS configuration.
ADD etc/litefs.yml /etc/litefs.yml

# Ensure our mount & data directories exists before mounting with LiteFS.
RUN mkdir -p /var/lib/hermes /mnt/litefs

ENTRYPOINT ["litefs"]
