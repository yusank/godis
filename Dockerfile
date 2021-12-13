FROM golang:1.17-alpine as builder

WORKDIR /app/godis
COPY . .

# args
ARG build_tags
RUN CGO_ENABLED=0 go build -tags "${build_tags}" -o godis cmd/server/main.go

FROM scratch

WORKDIR /app/godis

COPY --from=builder /app/godis .

EXPOSE 7379

CMD ["./godis"]
