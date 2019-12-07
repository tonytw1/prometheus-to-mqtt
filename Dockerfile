FROM golang:1.11

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go get github.com/stretchr/testify/
RUN go test -v ./...
run go build -v

CMD ["app"]
