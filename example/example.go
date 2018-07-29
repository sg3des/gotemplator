//go:generate gotemplator ./templates

package main

import (
	"net/http"

	"github.com/sg3des/gotemplator/example/templates"
)

func main() {
	http.HandleFunc("/", GoTemplator)
	http.HandleFunc("/writer/", GoTemplatorWriter)
	http.ListenAndServe("127.0.0.1:8090", nil)
}

func GoTemplator(w http.ResponseWriter, r *http.Request) {
	w.Write(templates.Templator([]string{"User0", "User1", "User1"}))

}

func GoTemplatorWriter(w http.ResponseWriter, r *http.Request) {
	templates.TemplatorWriter(w, []string{"User0", "User1", "User1"})
}
