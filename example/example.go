package main

import (
	"fmt"
	"local/gotemplator/example/templates"
	"net/http"
)

//go:generate gotemplator ./templates
func main() {
	fmt.Println("start")

	http.Handle("/", http.HandlerFunc(Route))
	http.ListenAndServe("127.0.0.1:8090", nil)

	fmt.Println("end")
}

func Route(w http.ResponseWriter, r *http.Request) {
	n, err := w.Write(templates.Index("World", map[int]string{34: "Val1", 35345: "Val2"}))
	fmt.Println(n, err)
}
