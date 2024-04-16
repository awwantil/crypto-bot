FROM golang:alpine AS builder

ENV CGO_ENABLED 0

ENV GOOS linux

WORKDIR /build

ADD ./go.mod .

COPY . .

RUN go build -trimpath -o backend ./backend-service/backend.go

#COPY ./.env ./build/.env

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /build

COPY --from=builder /build/backend /build/backend

#COPY --from=builder ./build/.env /build/.env

EXPOSE 8020

CMD ["/build/backend"]
