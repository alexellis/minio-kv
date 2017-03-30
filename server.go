package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"io/ioutil"

	"bytes"

	"strings"

	"github.com/gorilla/mux"
	minio "github.com/minio/minio-go"
)

func getHandler(minioClient *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		object := vars["object"]

		objectReceived, getErr := minioClient.GetObject("tables", object)
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
			fmt.Printf("%s\n%s", object, string(data))
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

		n, putErr := minioClient.PutObject("tables", object, reader, "application/text")
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

func get(key string) string {
	env := os.Getenv(key)
	if len(env) == 0 {
		path := os.Getenv(key + ".secret")
		data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalln(path, err)
		}

		env = strings.Trim(strings.SplitAfter(string(data), "\n")[0], "\n")
	}
	return env
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
			break;
		} else {
			log.Println(err)
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

func main() {
	ssl := false
	secret := get("secret")
	access := get("access")
	host := os.Getenv("host")

	fmt.Printf("secret='%s',access='%s'\n", secret, access)
	minioClient, err := connect(ssl, secret, access, host)
        if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Handle("/put/{object:[a-zA-Z0-9.-_]+}", putHandler(minioClient))
	r.Handle("/get/{object:[a-zA-Z0-9.-_]+}", getHandler(minioClient))

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}
	log.Fatal(s.ListenAndServe())
}
