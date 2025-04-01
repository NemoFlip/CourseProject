# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /auth

COPY ../go.mod ../go.sum ./

RUN go mod download

COPY auth_service/ auth_service/

COPY docs/ docs/

RUN go build -o ./auth_main ./auth_service/cmd/app/main.go

# Stage 2: final image
FROM alpine:latest

COPY --from=builder auth/auth_main /bin/auth_service

CMD ["/bin/auth_service"]