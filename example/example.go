package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	"github.com/sg3des/gotemplator/example/templates"
)

var (
	t   *template.Template
	tpl = `<!DOCTYPE html>
<html>
	<head>
		<title></title>
	</head>
	<body>
		<h1>Hello, {{.Name}}</h1>
		{{range $key, $val := .Data}}
			<div><b>{{$key}}</b>:<i>{{$val}}</i></div>
		{{end}}
	</body>
</html>`
)

//go:generate gotemplator ./templates
func main() {
	log.SetFlags(log.Lshortfile)
	var err error
	t, err = template.New("webpage").Parse(tpl)
	if err != nil {
		panic(err)
	}
	fmt.Println("start")

	http.Handle("/GoTemplator", http.HandlerFunc(GoTemplator))
	http.Handle("/Native", http.HandlerFunc(Native))
	http.ListenAndServe("127.0.0.1:8090", nil)

	fmt.Println("end")
}

func GoTemplator(w http.ResponseWriter, r *http.Request) {
	templates.Writer(w, "World", map[int]string{34: "Val1", rand.Int(): "Val2"})
	// log.Println("GoTemplator")
	// w.Write()
	// fmt.Println(err)
}

type Data struct {
	Name string
	Data map[int]string
}

func Native(w http.ResponseWriter, r *http.Request) {
	// log.Println("Native")
	t.Execute(w, Data{
		Name: "World",
		Data: map[int]string{34: "Val1", rand.Int(): "Val2"},
	})
	// fmt.Println(err)
}
