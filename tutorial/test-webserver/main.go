package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
)

var token, _ = randToken(3)

func main() {
	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/pod/", PodServer)
	http.HandleFunc("/health", HealthServer)
	fmt.Println("go app is up and running on host: ", getEnv("HOSTNAME", "not-set"))
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handling main request with color: ", token)
	fmt.Fprintf(w, "<h1 style=background-color:#%s>Hello %s</h1><h2>ID: %s</h2>", token, getEnv("NAME", "golang"), token)
}

func PodServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handling pod request")
	fmt.Fprintf(w, "%s", getEnv("HOSTNAME", "not-set"))
}

func HealthServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func randToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
