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

2) if the line like `||=Header(somevariable)` - then should call another section or any function returned `[]byte`, it's converted to `w.Write(Header(somevariable))`

3) code in `{{=var}}` - is short print variable, converted to `fmt.Fprintf(w, "%v", var)`

4) simple ternary operator - `{{?condition?then:else}}`, ex: `class='{{?len(val)>10?"long":"short"}}'`

everything else is html - which moved to go file how `w.Write([]byte(string))`

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
			|| for i, val := range users { 
				<li>{{?i==0?"selected":i}} {{=val}}</li>
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
```

If in template arguments exist `w io.Writer`, then html code will be written to it:

```
|| template IndexWriter(w io.Writer, users []string)
	<!DOCTYPE html>
	<html>
		<head>
			<title></title>
		</head>
		<body>
			<ul>
			|| for i, val := range users { 
				<li>{{?i==0?"selected":i}} {{=val}}</li>
				// {{=key}} - this is comment and will not be in the .go file
			|| } 
			</ul>
		</body>
	</html>
|| end
```

```go
func IndexWriter(w io.Writer, users []string) {
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
}
```

Templator support ternary operator:

`{{?if>10?"something"}}` converted to:

```go
if > 10 {
	fmt.Fprintf(w, "%v", "something")
}
```

and with else condition - `{{?i>10?"something":"if condition false"}}`:

```go
if i > 10 {
	fmt.Fprintf(w, "%v", "something")
} else {
	fmt.Fprintf(w, "%v", "if condition false")
}
```



### PERFOMANCE

GoTemplator slower than [Hero](http://github.com/shiyanhui/hero/) as uses `fmt.Fprintf()` for write user varables and internal creation of `bytes.Buffer`.

	BenchmarkTemplator-8         	  500000	      4157 ns/op	     928 B/op	       7 allocs/op
	BenchmarkTemplatorWriter-8   	 1000000	      1000 ns/op	     531 B/op	       3 allocs/op
	BenchmarkHero-8              	 2000000	       553 ns/op	     453 B/op	       0 allocs/op
