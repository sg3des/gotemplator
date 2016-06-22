package main

import (
	"errors"
	"fmt"
	"testing"
)

func init() {
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

	out := Print(in)

	if out != must {
		t.Error(errors.New("print html code is different"))
	}
}

func TestGoPrint_print(t *testing.T) {
	in := `{{=name}}`
	must := `fmt.Fprintf(_W, "%v", name)`

	out := GoPrint(in)
	if out != must {
		t.Error(errors.New("print go code is different! received " + out + " must be " + must))
	}

}

func TestGoPrint_raw(t *testing.T) {
	in := `{{strings.Join(aname,"\n")}}`
	must := `strings.Join(aname,"\n")`

	out := GoPrint(in)
	if out != must {
		t.Error(errors.New("print go code is different! received " + out + " must be " + must))
	}
}

func TestScan(t *testing.T) {
	line := `<div>Hello, {{=name}}!</div>`
	result := Scan(line)
	if len(result) != 3 {
		t.Error(errors.New("parse line failed"))
	}

	line = `// {{=name}}this is comment`
	result = Scan(line)
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
