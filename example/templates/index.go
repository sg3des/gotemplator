package templates

import (
	"bytes"
	"fmt"
	"io"
)

func Index(users []string) []byte {
	w := new(bytes.Buffer)
	w.Write([]byte("\t<!DOCTYPE html>\n\t<html>\n\t\t<head>\n\t\t\t<title></title>\n\t\t</head>\n\t\t<body>\n\t\t\t<ul>"))
	for i, val := range users {
		w.Write([]byte("\t\t\t\t<li>"))
		if i == 0 {
			fmt.Fprintf(w, "%v", "selected")
		} else {
			fmt.Fprintf(w, "%v", i)
		}
		fmt.Fprintf(w, " %v</li>", val)
	}
	w.Write([]byte("\t\t\t</ul>\n\t\t</body>\n\t</html>"))
	return w.Bytes()
}

func Countries(w io.Writer) {
	w.Write([]byte("\t<div>this template with writer</div>"))
	fmt.Fprintf(w, "\t<div>%v</div>", "one line print")
}
