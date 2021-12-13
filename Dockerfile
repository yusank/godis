# syntax=docker/dockerfile:1
FROM golang:1.17-alpine as builder

WORKDIR /app/godis

COPY . .
RUN CGO_ENABLED=0 go build -o godis cmd/server/main.go

FROM scratch

WORKDIR /app/godis

COPY --from=builder /app/godis .

EXPOSE 7379

CMD ["./godis"]
