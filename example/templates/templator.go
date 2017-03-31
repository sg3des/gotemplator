package templates

import (
	"bytes"
	"fmt"
	"io"
)

func Templator(users []string) []byte {
	w := new(bytes.Buffer)
	w.Write([]byte("\t<!DOCTYPE html>\n\t<html>\n\t\t<head>\n\t\t\t<title></title>\n\t\t</head>\n\t\t<body>\n\t\t\t<ul>"))
	for _, val := range users {
		fmt.Fprintf(w, "\t\t\t\t<li>%v</li>", val)
	}
	w.Write([]byte("\t\t\t</ul>\n\t\t</body>\n\t</html>"))
	return w.Bytes()
}

func TemplatorWriter(w io.Writer, users []string) {
	w.Write([]byte("\t<!DOCTYPE html>\n\t<html>\n\t\t<head>\n\t\t\t<title></title>\n\t\t</head>\n\t\t<body>\n\t\t\t<ul>"))
	for _, val := range users {
		fmt.Fprintf(w, "\t\t\t\t<li>%v</li>", val)
	}
	w.Write([]byte("\t\t\t</ul>\n\t\t</body>\n\t</html>"))
}
