package templates

import (
	"bytes"

	"github.com/shiyanhui/hero"
)

func Hero(users []string, buffer *bytes.Buffer) {
	buffer.WriteString(`

<!DOCTYPE html>
<html>
	<head>
		<title></title>
	</head>
	<body>
		<ul>
		`)
	for _, val := range users {
		buffer.WriteString(`
			<li>`)
		hero.EscapeHTML(val, buffer)
		buffer.WriteString(`</li>
		`)
	}
	buffer.WriteString(`
		</ul>
	</body>
</html>`)

}
