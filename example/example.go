package main

import (
	"fmt"
	"net/http"

	"github.com/sg3des/templator"
)

//go:generate gotemplator ./
func main() {
	fmt.Println("start")

	http.Handle("/", http.HandlerFunc(Route))
	http.ListenAndServe("127.0.0.1:8080", nil)

	fmt.Println("end")
}

func Route(w http.ResponseWriter, r *http.Request) {
	Index(w, []byte("World"), map[string]string{"Key1": "Val1", "Key2": "Val2"})
}
