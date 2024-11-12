FROM golang:alpine

WORKDIR /go-rest-api
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./bin/api ./cmd/api \
    && go build -o ./bin/migrate ./cmd/migrate \
    && ls -l ./bin

CMD ["/go-rest-api/bin/api"]

EXPOSE 8080
