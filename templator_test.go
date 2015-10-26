package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestPrint(t *testing.T) {
	in := `<html>`
	must := fmt.Sprintf(`w.Write([]byte("%s"))`, in)

	out := Print(in)

	if out != must {
		t.Error(errors.New("print html code is different"))
	}
}

func TestGoPrint_print(t *testing.T) {
	in := `{{=name}}`
	must := `fmt.Fprintf(w, "%s", name)`

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
}

func TestParse(t *testing.T) {
	gocode, err := Parse("./example/index.gtm")
	if err != nil {
		t.Error(err)
	}
	if len(gocode) == 0 {
		t.Error(errors.New("length of returned code is zero"))
	}
}
