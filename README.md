# minio-kv

minio-kv creates a simple key-value store from [Minio](https://min.io) or AWS S3. It adds a simple HTTP API for storing JSON documents and binary blobs with streaming capabilities.

## Usage

* Putting a JSON object:

    URI: `/put/{object:[a-zA-Z0-9.-_]+}`

    Request Body: contents of object

* Getting a JSON object:

    URI: `/get/{object:[a-zA-Z0-9.-_]+}`

    Response body: contents of object if found

* Putting a binary object with streaming:

    URI: `/put-streaming-blob/{object:[a-zA-Z0-9.-_]+}`

    Request Body: contents of object

* Getting a binary object with streaming:

    URI: `/put-streaming-blob/{object:[a-zA-Z0-9.-_]+}`

    Response body: contents of object if found


* Putting a binary object without streaming:

    URI: `/put-blob/{object:[a-zA-Z0-9.-_]+}`

    Request Body: contents of object

* Getting a binary object without streaming:

    URI: `/put-blob/{object:[a-zA-Z0-9.-_]+}`

    Response body: contents of object if found

## Status

Suitable for testing and local use, authentication is not included in the HTTP API.

Backlog:

- [x] Add streaming modes
- [x] Migrate to Docker / Go 1.11
- [x] Use vendoring for dependencies
- [ ] Read secrets from `/var/openfaas/secrets` for use with OpenFaaS as a "serverless workload"
- [ ] Add shared secret Bearer token for HTTP API

## Testing

```sh
docker run -e MINIO_ACCESS_KEY=511e8d6c84ee65feda34efdcc5366281a22b6dfd -e MINIO_SECRET_KEY=d88b4816ac11f0ee5efcc5282d2fe9896162a1d6 --name minio -p 9000:9000 minio/minio server /tmp/
```

Run the app

```sh
go build && port=8081 host=127.0.0.1:9000 MINIO_ACCESS_KEY=511e8d6c84ee65feda34efdcc5366281a22b6dfd MINIO_SECRET_KEY=d88b4816ac11f0
ee5efcc5282d2fe9896162a1d6 ./minio-kv
```

Get/put a big binary file:

```
curl localhost:8081/put-blob-stream/test --data-binary @big-file.img
curl localhost:8081/get-blob-stream/test > big-file-out.img

diff big-file.img big-file-out.img
```

## License

MIT, Alex Ellis 2017-2019
