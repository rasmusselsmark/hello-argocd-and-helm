package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, Helm and ArgoCD world!")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is listening on :80...")
	http.ListenAndServe(":80", nil)
}
