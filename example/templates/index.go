package templates

import (
	"bytes"
	"fmt"
)

func Index(users []string) []byte {
	_W := new(bytes.Buffer)

	_W.WriteString("<!DOCTYPE html><html><head><title></title></head><body><ul>")
	for _, val := range users {
		_W.WriteString("<li>")
		fmt.Fprintf(_W, "%v", val)
		_W.WriteString("</li>")

	}

	_W.WriteString("</ul></body></html>")
	return _W.Bytes()
}
