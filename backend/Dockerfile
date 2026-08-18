FROM golang:1.25.1-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/detective-game-api ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app
COPY --from=builder --chown=app:app /out/detective-game-api ./detective-game-api

USER app

EXPOSE 8080

ENTRYPOINT ["./detective-game-api"]
