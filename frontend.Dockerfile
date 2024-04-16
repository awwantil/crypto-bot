FROM golang:alpine AS builder

ENV CGO_ENABLED 0

ENV GOOS linux

WORKDIR /build

ADD ./go.mod .

COPY . .

RUN go build -trimpath -o frontend ./frontend-service/frontend.go

#COPY ./.env ./build/.env

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /build

COPY --from=builder /build/frontend /build/frontend

#COPY --from=builder ./build/.env /build/.env

EXPOSE 8000

CMD ["/build/frontend"]
