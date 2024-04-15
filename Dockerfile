FROM golang:1.22.2-alpine3.19 as builder

WORKDIR /build

COPY . /build/

RUN CGO_ENABLED=0 GOOS=linux go build -o busofact .

FROM alpine

RUN adduser -S -D -H -h /app appuser

USER appuser

COPY . /app

COPY --from=builder /build/busofact /app/

WORKDIR /app

EXPOSE 8080

CMD [ "./busofact" ]