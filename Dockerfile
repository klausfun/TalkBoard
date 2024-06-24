FROM golang:1.22 AS builder

RUN go version
ENV GOPATH=/

COPY ./ ./

ENV STORAGE_TYPE=${STORAGE_TYPE}

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o talkBoard ./cmd/main.go

CMD ["./talkBoard", "-storage", "${STORAGE_TYPE}"]
