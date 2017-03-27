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
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Can't read object " + object)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Printf("'%s'", string(data))
		w.Write(data)
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
		}
		log.Printf("n=%d\n", n)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Put " + object + ". OK."))

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

func main() {
	ssl := false
	secret := get("secret")
	access := get("access")
	fmt.Printf("secret='%s',access='%s'\n", secret, access)
	var minioClient *minio.Client
	var err error
	for i := 0; i < 10; i++ {
		fmt.Printf("Connecting: %d/5\n", i)
		minioClient, err = minio.New(os.Getenv("host"),
			access,
			secret, ssl)

		if err != nil {
			log.Println(err)
		}
		time.Sleep(1 * time.Second)
	}
        if err != nil {
            log.Fatal("Cannot connect to S3")
            return
        }
	err = minioClient.MakeBucket("tables", "us-east-1")
	if err != nil {
		fmt.Println(err)
                return
	} else {
		fmt.Println("Successfully created bucket \"tables\".")
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
