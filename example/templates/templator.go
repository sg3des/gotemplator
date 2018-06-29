package templates

import (
	"bytes"
	"fmt"
	"io"
)

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
	for _, val := range users {
		fmt.Fprintf(w, `				<li>%v</li>`, val)
	}
	w.Write(html3)
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
		<body>
			<ul>`)
var html3 = []byte(`			</ul>
		</body>
	</html>`)
