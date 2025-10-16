# Build
FROM golang:1.22 as build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/retech-core ./cmd/api

# Runtime
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /out/retech-core /retech-core
COPY internal/docs /internal/docs
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/retech-core"]

