package templates

import (
	"bytes"
	"fmt"
	"io"
)

func Templator(users []string) []byte {
	w := new(bytes.Buffer)
	w.Write([]byte(`	<!DOCTYPE html>
	<html>
		<head>
			<title></title>
		</head>
		<body>
			<ul>`))
	for _, val := range users {
		fmt.Fprintf(w, `				<li width='100%%'>%v</li><i width='50%%'></i>`, val)
	}
	w.Write([]byte(`			</ul>
		</body>
	</html>`))
	return w.Bytes()
}

func TemplatorWriter(w io.Writer, users []string) {
	w.Write([]byte(`	<!DOCTYPE html>
	<html>
		<head>
			<title></title>
		</head>
		<body>
			<ul>`))
	for _, val := range users {
		fmt.Fprintf(w, `				<li>%v</li>`, val)
	}
	w.Write([]byte(`			</ul>
		</body>
	</html>`))
}
