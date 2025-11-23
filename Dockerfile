FROM golang:1.25 AS build

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN go build -o server ./cmd ;

FROM debian:stable-slim

WORKDIR /app
COPY --from=build /app/server .

CMD ["/app/server"]
