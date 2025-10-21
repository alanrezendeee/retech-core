# Build
FROM golang:1.24 as build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/retech-core ./cmd/api

# Runtime
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/retech-core /app/retech-core
COPY internal/docs /app/internal/docs
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/retech-core"]

