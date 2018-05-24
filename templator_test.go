package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/sg3des/gotemplator/example/templates"
)

var dir string

func init() {
	log.SetFlags(log.Lshortfile)
	// *verbose = true
	dir = "./example/templates/"
}

func TestGenerate(t *testing.T) {
	gtms, err := getFiles(dir, ".gtm") //if failed will be panic
	if err != nil {
		t.Error(err)
	}

	if len(gtms) == 0 {
		t.Error("templates not found")
	}
}

func TestPrint(t *testing.T) {
	in := []byte(`<html>`)
	must := []byte(fmt.Sprintf(`w.Write([]byte("%s"))`, in))

	out := printHTML(in)

	if !bytes.Equal(out, must) {
		t.Error(fmt.Errorf("print html code is different, recieved '%s', must be '%s'", out, must))
	}
}

func TestGoPrint_print(t *testing.T) {
	in := []byte(`name`)
	must := []byte("fmt.Fprintf(w, \"%v\", name)")

	out := printGocode(nil, in, nil)
	if !bytes.Equal(out, must) {
		t.Error(fmt.Errorf("print go code is different! received `%s` must be `%s`", out, must))
	}

}

func TestScan(t *testing.T) {
	line := []byte(`<div>Hello, {{=name}}!</div>`)
	must := [][]byte{[]byte(`fmt.Fprintf(w, "<div>Hello, %v!</div>", name)`)}

	var writerExist bool
	var htmlLines [][]byte
	result := Scan(line, &htmlLines, &writerExist)

	if len(result) != len(must) {
		t.Error("length of result should be ", len(must))
	}
	if !bytes.Equal(result[0], must[0]) {
		t.Errorf("parse line failed! received: `%s` must: `%s`", string(result[0]), string(must[0]))
	}

	line = []byte(`// {{=name}}this is comment`)
	result = Scan(line, &htmlLines, &writerExist)
	if len(result) != 0 {
		t.Error(errors.New("parse comment line failed"))
	}

}

func TestParse(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	err := Parse("./example/templates/templator.gtm", buf)
	if err != nil {
		t.Error(err)
	}
}

func TestTernary(t *testing.T) {
	line := []byte("<div>{{?x>10?long:short}}")
	must := [][]byte{
		[]byte(`w.Write([]byte("<div>"))`),
		[]byte(`if x>10 {`),
		[]byte(`	fmt.Fprintf(w, "%v", long)`),
		[]byte(`} else {`),
		[]byte(`	fmt.Fprintf(w, "%v", short)`),
		[]byte(`}`),
	}

	var writerExist bool
	var htmlLines [][]byte
	out := Scan(line, &htmlLines, &writerExist)

	if len(out) != len(must) {
		t.Errorf("length shoud be equal %d != %d", len(out), len(must))
	}
	for i, m := range must {
		o := out[i]
		if !bytes.Equal(o, m) {
			t.Errorf("not equal `%s` `%s`", string(o), string(m))
		}
	}
}

func BenchmarkTemplator(b *testing.B) {
	userlist := []string{
		"Alice",
		"Bob",
		"Tom",
	}
	for n := 0; n < b.N; n++ {
		templates.Templator(userlist)
	}
}

func BenchmarkTemplatorWriter(b *testing.B) {
	var w = new(bytes.Buffer)
	userlist := []string{
		"Alice",
		"Bob",
		"Tom",
	}
	for n := 0; n < b.N; n++ {
		templates.TemplatorWriter(w, userlist)
	}
}

func BenchmarkHero(b *testing.B) {
	var w = new(bytes.Buffer)
	userlist := []string{
		"Alice",
		"Bob",
		"Tom",
	}
	for n := 0; n < b.N; n++ {
		templates.Hero(userlist, w)
	}
}
