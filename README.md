# minio-db

A thin layer on top of Minio-db to store JSON objects with GET/PUT over HTTP.

### Usage:

Putting an object:

URI: /put/{object:[a-zA-Z0-9.-_]+}

Request Body: contents of object

Getting an object:

URI: /get/{object:[a-zA-Z0-9.-_]+}

Response body: contents of object if found

### Building:

This can be built with Docker or Golang and can accept a Docker swarm secret or environmental variable for configuration of Minio secret/access key.

