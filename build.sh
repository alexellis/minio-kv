#!/bin/sh

docker build --build-arg http_proxy=$http_proxy --build-arg https_proxy=$https_proxy -t alexellis2/minio-kv:0.1.1 .

