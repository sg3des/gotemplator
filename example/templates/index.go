package templates

import "fmt"

func Index(name string, data map[int]string) {
	_W.WriteString("<!DOCTYPE html>\n")
	_W.WriteString("<html>\n")
	_W.WriteString("\t<head>\n")
	_W.WriteString("\t\t<title></title>\n")
	_W.WriteString("\t</head>\n")
	_W.WriteString("\t<body>\n")
	_W.WriteString("\t\t<h1>Hello, ")
	fmt.Fprintf(_W, "%v", name)
	_W.WriteString("</h1>\n")
	for key, val := range data {
		_W.WriteString("\t\t\t<div><b>")
		fmt.Fprintf(_W, "%v", key)
		_W.WriteString("</b>:<i>")
		fmt.Fprintf(_W, "%v", val)
		_W.WriteString("</i></div>\n")
	}
	_W.WriteString("\t</body>\n")
	_W.WriteString("</html>\n")
	return _W.Bytes()
}
