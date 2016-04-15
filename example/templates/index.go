package templates

import (
	"bytes"
	"fmt"
)

func Index(name string, data map[int]string) []byte {
	_W := bytes.NewBuffer([]byte{})

	_W.WriteString("\t<!DOCTYPE html>\n")
	_W.WriteString("\t<html>\n")
	_W.WriteString("\t\t<head>\n")
	_W.WriteString("\t\t\t<title></title>\n")
	_W.WriteString("\t\t</head>\n")
	_W.WriteString("\t\t<body>\n")
	_W.WriteString("\t\t\t<h1>Hello, ")
	fmt.Fprintf(_W, "%v", name)
	_W.WriteString("</h1>\n")
	for key, val := range data {
		_W.WriteString("\t\t\t\t<div><b>")
		fmt.Fprintf(_W, "%v", key)
		_W.WriteString("</b>:<i>")
		fmt.Fprintf(_W, "%v", val)
		_W.WriteString("</i></div>\n")
	}
	_W.WriteString("\t\t</body>\n")
	_W.WriteString("\t</html>\n")
	return _W.Bytes()
}
