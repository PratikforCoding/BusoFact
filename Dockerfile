FROM golang:1.22.2-alpine3.19

WORKDIR /app

COPY . /app/

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /BusoFact

EXPOSE 8080

CMD [ "/BusoFact" ]
