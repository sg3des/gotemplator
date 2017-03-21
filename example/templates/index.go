package templates

import (
	"bytes"
	"fmt"
)

func Index(users []string) []byte {
	_W := new(bytes.Buffer)

	_W.WriteString("\t<!DOCTYPE html>\t<html>\t\t<head>\t\t\t<title></title>\t\t</head>\t\t<body>\t\t\t<ul>")
	for _, val := range users {
		_W.WriteString("\t\t\t\t<li>")
		fmt.Fprintf(_W, "%v", val)
		_W.WriteString("</li>")

	}

	_W.WriteString("\t\t\t</ul>\t\t</body>\t</html>")
	return _W.Bytes()
}
