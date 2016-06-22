# TEMPLATOR

Templator is alternative view on html templates for go - on my opinion, standard package "text/templates" offer awful syntax! 

Why syntax in the templates must be different from the GO syntax?

1. Coplex html templates become pretty hard to understand.
2. Functions\pipeline calls are also very poor, we have to write functions such as eq, more, less, et.c.

Templator generate **go** code from simply templates(*.gtm).

HTML template - it function, with full syntax GO. Template functions returned []byte HTML code.


## USAGE

- manually execute `gotemplator /path/to/dir`

OR use go generate:

- add to your project:	`//go:generate gotemplator /path/to/dir/` and execute `go generate`

## SYNTAX

Template syntax is hybrid of html and go, have only 3 rules:

1) if the line begins with `||` - then this is go code which is unchanged moved to go file

2) code in `{{ }}`(double curly braces) - is too go code

3) code in `{{=var}}` - is short print variable

everything else is html - which moved to go file how `_W.WriteString(string)`

example:
	
	|| template Index(name string, data map[string]string) {
	<!DOCTYPE html>
	<html>
	    <head>
	        <title></title>
	    </head>
	    <body>
	        <h1>Hello, {{=name}}</h1>
	        {{for key,val := range data { }}
	          <div><b>{{=key}}</b>:<i>{{=val}}</i></div>
	        {{ } }}
	    </body>
	</html>
	|| end

this is transform to:

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

