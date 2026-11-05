FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/rest ./cmd/rest \
 && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/cli ./cmd/cli \
 && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

FROM alpine:3.22
RUN apk add --no-cache sqlite tzdata \
 && addgroup -S -g 65532 abacus \
 && adduser -S -D -H -u 65532 -G abacus abacus \
 && mkdir /data && chown abacus:abacus /data
COPY --from=build /out/ /usr/local/bin/
USER abacus
VOLUME /data
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/rest"]
