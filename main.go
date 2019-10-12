package main

import (
	"crypto/subtle"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"io/ioutil"

	"bytes"

	"github.com/gorilla/mux"
	minio "github.com/minio/minio-go"
)

type bearerToken struct {
	token string
}

func getBlobHandler(minioClient *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		object := vars["object"]
		getOpts := minio.GetObjectOptions{}

		objectReceived, getErr := minioClient.GetObject("tables", object, getOpts)
		if getErr != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println("Can't get object " + object)
			return
		}

		data, readErr := ioutil.ReadAll(objectReceived)
		if readErr != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println("Can't read object " + object)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func getBlobStreamHandler(minioClient *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		object := vars["object"]
		getOpts := minio.GetObjectOptions{}

		objectReceived, getErr := minioClient.GetObject("tables", object, getOpts)
		if getErr != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println("Can't get object " + object)
			return
		}

		rr := ioutil.NopCloser(objectReceived)
		defer rr.Close()

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, rr)
	}
}

func putBlobHandler(minioClient *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		object := vars["object"]

		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		reader := bytes.NewReader(body)

		putOpts := minio.PutObjectOptions{
			ContentEncoding: "application/octet-stream",
		}

		n, putErr := minioClient.PutObject("tables", object, reader, int64(len(body)), putOpts)

		if putErr != nil {
			log.Println("Can't put object " + object)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Can't put " + object + ". Failed."))
		} else {
			log.Printf("Put %s. n=%d\n", object, n)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Put " + object + ". OK."))
		}
	}
}

func putBlobStreamHandler(minioClient *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		object := vars["object"]

		defer r.Body.Close()

		putOpts := minio.PutObjectOptions{
			ContentEncoding: "application/octet-stream",
		}

		n, putErr := minioClient.PutObject("tables", object, r.Body, -1, putOpts)

		if putErr != nil {
			log.Println("Can't put object " + object)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Can't put " + object + ". Failed."))
		} else {
			log.Printf("Put %s. n=%d\n", object, n)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Put " + object + ". OK."))
		}
	}
}

func getHandler(minioClient *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		object := vars["object"]

		getOpts := minio.GetObjectOptions{}

		objectReceived, getErr := minioClient.GetObject("tables", object, getOpts)
		if getErr != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println("Can't get object " + object)
			return
		}

		data, readErr := ioutil.ReadAll(objectReceived)
		if readErr != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println("Can't read object " + object)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Printf("Read: %s, bytes: %d\n", object, len(data))
			w.Write(data)
		}
	}
}

func putHandler(minioClient *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		object := vars["object"]

		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		reader := bytes.NewReader(body)

		putOpts := minio.PutObjectOptions{
			ContentEncoding: "application/json",
		}

		n, putErr := minioClient.PutObject("tables", object, reader, int64(len(body)), putOpts)

		if putErr != nil {
			log.Println("Can't put object " + object)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Can't put " + object + ". Failed."))
		} else {
			log.Printf("Put %s. n=%d\n", object, n)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Put " + object + ". OK."))
		}
	}
}

func connect(ssl bool, secret string, access string, host string) (*minio.Client, error) {
	maxAttempts := 30
	connected := false
	var err error
	var minioClient *minio.Client

	for i := 1; i <= maxAttempts; i++ {
		fmt.Printf("Connecting: %d/%d\n", i, maxAttempts)
		minioClient, err = minio.New(host, access, secret, ssl)

		if err == nil {
			connected = true
			break
		} else {
			log.Printf("Error: %s\n", err)
			time.Sleep(1 * time.Second)
		}
	}

	if connected == false && err != nil {
		log.Fatal("Cannot connect to S3")
		return nil, err
	}

	exists, err := minioClient.BucketExists("tables")
	if err == nil && exists == false {
		err = minioClient.MakeBucket("tables", "us-east-1")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Successfully created bucket \"tables\".")
		}
	}

	return minioClient, err
}

func (t *bearerToken) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := ""
		if tokenSplit := strings.Split(r.Header.Get("Authorization"), "Bearer"); len(tokenSplit) > 1 {
			token = strings.TrimSpace(tokenSplit[1])
		}

		if subtle.ConstantTimeCompare([]byte(token), []byte(t.token)) == 1 {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}

func main() {
	ssl := false
	secret := os.Getenv("MINIO_SECRET_KEY")
	access := os.Getenv("MINIO_ACCESS_KEY")
	host := os.Getenv("host")
	port := "8080"
	if val, ok := os.LookupEnv("port"); ok {
		port = strings.TrimSpace(val)
	}

	bearerToken := bearerToken{}
	flag.StringVar(&bearerToken.token, "token", "", "pass token for client to authenticate")
	flag.Parse()

	minioClient, err := connect(ssl, secret, access, host)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	if len(bearerToken.token) > 0 {
		r.Use(bearerToken.authenticate)
	}

	r.Handle("/put/{object:[a-zA-Z0-9.-_]+}", putHandler(minioClient))
	r.Handle("/get/{object:[a-zA-Z0-9.-_]+}", getHandler(minioClient))

	r.Handle("/put-blob/{object:[a-zA-Z0-9.-_]+}", putBlobHandler(minioClient))
	r.Handle("/get-blob/{object:[a-zA-Z0-9.-_]+}", getBlobHandler(minioClient))
	r.Handle("/put-blob-stream/{object:[a-zA-Z0-9.-_]+}", putBlobStreamHandler(minioClient))
	r.Handle("/get-blob-stream/{object:[a-zA-Z0-9.-_]+}", getBlobStreamHandler(minioClient))

	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		ReadTimeout:    5 * time.Minute,
		WriteTimeout:   5 * time.Minute,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}
	log.Fatal(s.ListenAndServe())
}
