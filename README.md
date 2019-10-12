# minio-kvp

minio-kvp is a thin layer to store JSON objects and binary blobs in a Minio object storage server with a GET/PUT interface over HTTP. For more on Minio, checkout https://minio.io

### Usage:

Putting an object:

URI: /put/{object:[a-zA-Z0-9.-_]+}

Request Body: contents of object

Getting an object:

URI: /get/{object:[a-zA-Z0-9.-_]+}

Response body: contents of object if found

### Building:

This can be built with Docker or Golang and can accept a Docker swarm secret or environment variable for configuration of Minio secret/access key.

### Testing

```sh
docker run -e MINIO_ACCESS_KEY=511e8d6c84ee65feda34efdcc5366281a22b6dfd -e MINIO_SECRET_KEY=d88b4816ac11f0ee5efcc5282d2fe9896162a1d6 --name minio -p 9000:9000 minio/minio server /tmp/
```

Run the app

```sh
go build && port=8081 host=127.0.0.1:9000 MINIO_ACCESS_KEY=511e8d6c84ee65feda34efdcc5366281a22b6dfd MINIO_SECRET_KEY=d88b4816ac11f0
ee5efcc5282d2fe9896162a1d6 ./minio-kvp
```

Get/put a big binary file:

```
curl localhost:8081/put-blob-stream/test --data-binary @big-file.img
curl localhost:8081/get-blob-stream/test > big-file-out.img

diff big-file.img big-file-out.img
```
