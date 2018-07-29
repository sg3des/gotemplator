package templates

import (
	"bytes"
	"fmt"
	"io"
)

func somefunc(w io.Writer, s string) {
	fmt.Fprintln(w, s)
}

func Templator(users []string) []byte {
	w := new(bytes.Buffer)
	w.Write(html0)
	for _, val := range users {
		fmt.Fprintf(w, `				<li width='100%%'>%v</li><i width='50%%'></i>`, val)
	}
	w.Write(html1)
	return w.Bytes()
}

func TemplatorWriter(w io.Writer, users []string) {
	w.Write(html2)
	somefunc(w, "text")
	w.Write(html3)
	for _, val := range users {
		fmt.Fprintf(w, `				<li>%v</li>`, val)
	}
	w.Write(html4)
}

var html0 = []byte(`	<!DOCTYPE html>
	<html>
		<head>
			<title></title>
		</head>
		<body>
			<ul>`)
var html1 = []byte(`			</ul>
		</body>
	</html>`)
var html2 = []byte(`	<!DOCTYPE html>
	<html>
		<head>
			<title></title>
		</head>
		<body>`)
var html3 = []byte(`			<ul>`)
var html4 = []byte(`			</ul>
		</body>
	</html>`)
