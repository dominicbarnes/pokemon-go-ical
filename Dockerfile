# Build the application from source
FROM golang:1.21 AS build-stage

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/ ./cmd/...


# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/bin/* /bin/

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/bin/service"]