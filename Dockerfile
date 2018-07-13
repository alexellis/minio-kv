FROM golang:1.7.5

RUN go get github.com/minio/minio-go && \
    go get github.com/gorilla/mux && \
    go get github.com/gorilla/context

RUN mkdir -p /go/src/github.com/alexellis/minio-db
WORKDIR /go/src/github.com/alexellis/minio-db

COPY server.go .
RUN go build -o server

EXPOSE 8080

CMD ["./server"]
