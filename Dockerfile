FROM golang:1.11

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
run go build -v

CMD ["app"]
