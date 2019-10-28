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
- [ ] Add Kubernetes `Deployment` and `Service`
- [ ] Read secrets from either env-var or from files mounted at path
- [ ] Add auth for all HTTP APIs using shared secret Bearer token

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

## Running on Kubernetes

**Create MINIO_ACCESS_KEY, MINIO_SECRET_KEY environment variables**
```
MINIO_ACCESS_KEY=$(head -c 12 /dev/urandom | shasum | cut -d  ' ' -f1) \
MINIO_SECRET_KEY=$(head -c 12 /dev/urandom | shasum | cut -d  ' ' -f1) \
TOKEN=$(head -c 12 /dev/urandom | shasum | cut -d  ' ' -f1)
```

**Create secrets out of the above env vars**
```
kubectl create secret generic minio-kv \
--from-literal=minio-access-key=$MINIO_ACCESS_KEY \
--from-literal=minio-secret-key=$MINIO_SECRET_KEY \
--from-literal=token=$TOKEN
```

**(If Helm is not installed)**
Follow the instructions [here](https://github.com/openfaas/faas-netes/blob/master/HELM.md).

**Install Minio with Helm**
```
helm install --name minio --namespace default \
   --set accessKey="$MINIO_ACCESS_KEY" \
   --set secretKey="$MINIO_SECRET_KEY" \
   --set replicas=1 \
   --set persistence.enabled=false \
   --set service.port=9000 \
   --set service.type=ClusterIP \
stable/minio
```

**Create the kubernetes deployment and service**
```
kubectl apply -f yaml/
```

## License

MIT, Alex Ellis 2017-2019
