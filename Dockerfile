FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/abacus ./cmd/rest

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/abacus /abacus
VOLUME /data
EXPOSE 8080
ENTRYPOINT ["/abacus"]
