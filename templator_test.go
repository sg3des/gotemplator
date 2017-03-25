package main

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/sg3des/gotemplator/example/templates"
)

func init() {
	log.SetFlags(log.Lshortfile)
	// *verbose = true
	dir = "./example/templates/"
}

func TestGenerate(t *testing.T) {
	gtms, err := Generate(dir) //if failed will be panic
	if err != nil {
		t.Error(err)
	}

	if len(gtms) == 0 {
		t.Error("templates not found")
	}
}

func TestPrint(t *testing.T) {
	in := `<html>`
	must := fmt.Sprintf(`_W.WriteString("%s")`, in)

	out := PrintHTML(in)

	if out != must {
		t.Error(errors.New("print html code is different"))
	}
}

func TestGoPrint_print(t *testing.T) {
	in := `{{=name}}`
	must := "fmt.Fprintf(_W, \"%v\", name)"

	out := GoPrint(in)
	if out != must {
		t.Error(fmt.Errorf("print go code is different! received `%s` must be `%s`", out, must))
	}

}

func TestGoPrint_raw(t *testing.T) {
	in := `{{strings.Join(aname,"\n")}}`
	must := `strings.Join(aname,"\n")`

	out := GoPrint(in)
	if out != must {
		t.Error(fmt.Errorf("print go code is different! received `%s` must be `%s`", out, must))
	}
}

func TestScan(t *testing.T) {
	line := `<div>Hello, {{=name}}!</div>`
	must := []string{`_W.WriteString("<div>Hello, ")`,
		`fmt.Fprintf(_W, "%v", name)`,
		`_W.WriteString("!</div>")`,
	}

	var htmlLines []string
	result := Scan(line, &htmlLines)

	if len(result) != len(must) {
		t.Error("length of result should be ", len(must))
	}
	if result[0] != must[0] {
		t.Error("parse line failed! received: `%s` must: `%s`", result[0], must[0])
	}
	if result[1] != must[1] {
		t.Errorf("parse line failed! received: `%s` must: `%s`", result[1], must[1])
	}
	if result[2] != must[2] {
		t.Errorf("parse line failed! received: `%s` must: `%s`", result[2], must[2])
	}

	line = `// {{=name}}this is comment`
	result = Scan(line, &htmlLines)
	if len(result) != 0 {
		t.Error(errors.New("parse comment line failed"))
	}

}

func TestParse(t *testing.T) {
	gocode, err := Parse("./example/templates/index.gtm")
	if err != nil {
		t.Error(err)
	}
	if len(gocode) == 0 {
		t.Error(errors.New("length of returned code is zero"))
	}
}

func TestTernary(t *testing.T) {
	line := "<div>{{?x>10?long:short}}"
	must := []string{`_W.WriteString("<div>")`,
		`if x>10 {
fmt.Fprintf(_W, "%v", long)
} else {
fmt.Fprintf(_W,"%v",short)
}`}

	var htmlLines []string
	gocode := Scan(line, &htmlLines)

	if len(gocode) != len(must) {
		t.Errorf("length shoud be equal", len(gocode), len(must))
	}
	if gocode[0] != must[0] {
		t.Errorf("not equal `%s` `%s`", gocode[0], must[0])
	}
	if gocode[1] != must[1] {
		t.Errorf("not equal recieved:\n`%s`\nmust:\n`%s`", gocode[1], must[1])
	}
}

func BenchmarkGoTemplator(b *testing.B) {
	// var w = new(bytes.Buffer)
	userlist := []string{
		"Alice",
		"Bob",
		"Tom",
	}
	for n := 0; n < b.N; n++ {
		templates.Index(userlist)
	}
}
