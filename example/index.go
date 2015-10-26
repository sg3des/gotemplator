package main

import (
	"fmt"
	"io"
)

func Index(w io.Writer, name []byte, data map[string]string) {
	w.Write([]byte("<!DOCTYPE html>"))
	w.Write([]byte("<html>"))
	w.Write([]byte("    <head>"))
	w.Write([]byte("        <title></title>"))
	w.Write([]byte("    </head>"))
	w.Write([]byte("    <body>"))
	w.Write([]byte("        <h1>Hello, "))
	fmt.Fprintf(w, "%s", name)
	w.Write([]byte("</h1>"))
	w.Write([]byte("        "))
	for key, val := range data {
		w.Write([]byte(""))
		w.Write([]byte("          <div><b>"))
		fmt.Fprintf(w, "%s", key)
		w.Write([]byte("</b>:<i>"))
		fmt.Fprintf(w, "%s", val)
		w.Write([]byte("</i></div>"))
		w.Write([]byte("        "))
	}
	w.Write([]byte(""))
	w.Write([]byte("    </body>"))
	w.Write([]byte("</html>"))
}
