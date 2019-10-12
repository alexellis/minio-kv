FROM golang:1.11 as build

WORKDIR github.com/alexellis/minio-kvp

COPY .git               .git
COPY vendor             vendor
COPY server.go            .

ARG GIT_COMMIT
ARG VERSION

RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.GitCommit=${GIT_COMMIT} -X main.Version=${VERSION}" -a -installsuffix cgo -o /usr/bin/minio-kvp

FROM alpine:3.10
RUN apk add --force-refresh ca-certificates

# Add non-root user
RUN addgroup -S app && adduser -S -g app app \
  && mkdir -p /home/app || : \
  && chown -R app /home/app

RUN touch /tmp/.lock

COPY --from=build /usr/bin/minio-kvp /usr/bin/
WORKDIR /home/app

USER app
EXPOSE 8080

CMD ["minio-kvp"]
