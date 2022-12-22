FROM golang:1.18 as builder
WORKDIR /app/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo -o app ./cmd/server

FROM scratch
COPY --from=builder /app/app /app
ENTRYPOINT ["/app"]
