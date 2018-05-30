package main

import (
	"net/http"

	"github.com/sg3des/gotemplator/example/templates"
)

//go:generate gotemplator ./templates
func main() {
	http.Handle("/", http.HandlerFunc(GoTemplator))
	http.ListenAndServe("127.0.0.1:8090", nil)
}

func GoTemplator(w http.ResponseWriter, r *http.Request) {
	w.Write(templates.Templator([]string{"User0", "User1", "User1"}))
}
