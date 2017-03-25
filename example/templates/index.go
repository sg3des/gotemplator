package templates

import (
	"bytes"
	"fmt"
)

func Index(users []string) []byte {
	_W := new(bytes.Buffer)

	_W.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<title></title>\n</head>\n<body>\n<ul>")
	for i, val := range users {
		_W.WriteString("<li>")
		if i == 0 {
			fmt.Fprintf(_W, "%v", "selected")
		} else {
			fmt.Fprintf(_W, "%v", i)
		}

		fmt.Fprintf(_W, "%v", val)
		_W.WriteString("</li>")

	}
	_W.WriteString("</ul>\n</body>\n</html>")
	return _W.Bytes()
}
