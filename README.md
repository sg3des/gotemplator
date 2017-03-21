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

Template syntax is hybrid of html and go, have only 4 rules:

1) if the line begins with `||` - then this is go code which is unchanged moved to go file

2) code in `{{ }}`(double curly braces) - is too go code

3) code in `{{=var}}` - is short print variable

4) if the line like `||=Header(somevariable)` - then should call another section or any function returned `[]byte`

everything else is html - which moved to go file how `_W.WriteString(string)`

example:
	
```html
|| template Index(users []string) 
	<!DOCTYPE html>
	<html>
		<head>
			<title></title>
		</head>
		<body>
			<ul>
			|| for _, val := range users { 
				<li>{{=val}}</li>
				// {{=key}} - this is comment and will not be in the .go file
			|| } 
			</ul>
		</body>
	</html>
|| end
```

this is transform to:

```go
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
```



### PERFOMANCE

GoTemplator slower than [Hero](http://github.com/shiyanhui/hero/) as uses `fmt.Fprintf()` for write user varables and internal creation of `bytes.Buffer`.

	BenchmarkGoTemplator-8  2000000	  767 ns/op	  400 B/op	  6 allocs/op
	BenchmarkHero-8         3000000	  334 ns/op	  698 B/op	  0 allocs/op
