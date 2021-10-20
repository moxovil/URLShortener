FROM golang:latest

RUN mkdir /app
COPY ./ /app

WORKDIR /app

#install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go build -o urllist-app ./cmd/main.go

CMD ["/app/urllist-app"]