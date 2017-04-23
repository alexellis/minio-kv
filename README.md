# minio-db

Minio-DB is a thin layer to store JSON objects and binary blobs in a Minio object storage server with a GET/PUT interface over HTTP. For more on Minio, checkout https://minio.io

### Usage:

Putting an object:

URI: /put/{object:[a-zA-Z0-9.-_]+}

Request Body: contents of object

Getting an object:

URI: /get/{object:[a-zA-Z0-9.-_]+}

Response body: contents of object if found

### Building:

This can be built with Docker or Golang and can accept a Docker swarm secret or environmental variable for configuration of Minio secret/access key.

