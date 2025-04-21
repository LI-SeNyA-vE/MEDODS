FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
#ENV CGO_ENABLED=0 \
#    GOOS=linux
RUN go build -o server cmd/server/main.go

FROM scratch
COPY --from=builder /app/server /server
EXPOSE 8901
ENTRYPOINT ["/server"]