//Templator is  alternative view on html templates for go!
//is the generator go code returned html

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sg3des/argum"

	"golang.org/x/tools/imports"
)

var args struct {
	Verbose bool   `argum:"-v,--verbose" help:"verbose mode"`
	Ext     string `argum:"-e,--ext" help:"file extension of templates" default:".gtm"`
	Path    string `argum:"pos,req" help:"path to directory with templates or file template"`
}

func init() {
	log.SetFlags(log.Lshortfile)

	argum.Version = "1.1.3.170329"
	argum.MustParse(&args)
	// arg.MustParse(&args)
}

func main() {
	files, err := getFiles(args.Path, args.Ext)
	if err != nil {
		log.Fatalln(err)
	}

	for _, gtmfilename := range files {
		if args.Verbose {
			fmt.Println(gtmfilename)
		}

		w := bytes.NewBuffer([]byte{})

		err := Parse(gtmfilename, w)
		if err != nil {
			log.Fatalln(err)
		}

		if args.Verbose {
			displayGeneratedCode(w.Bytes())
		}

		filedata := w.Bytes()
		filename := strings.TrimSuffix(gtmfilename, filepath.Ext(gtmfilename)) + ".go"
		opt := &imports.Options{AllErrors: true}
		importsData, err := imports.Process(filename, filedata, opt)
		if err != nil {
			displayGeneratedCode(filedata)
			log.Println(err)
			// importsData = filedata
			return
		}

		ioutil.WriteFile(filename, importsData, 0755)
	}
}

func getFiles(dir, ext string) ([]string, error) {
	fi, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return []string{dir}, nil
	}

	return filepath.Glob(path.Join(dir, "*"+args.Ext))
}

func displayGeneratedCode(data []byte) {
	lines := bytes.Split(data, []byte("\n"))
	for i, line := range lines {
		fmt.Printf("%d: %s\n", i, string(line))
	}
}

//Parse is parser for .gtm file and return go code
func Parse(filename string, w io.Writer) error {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}

	w.Write(addPackageLine(filename))

	var htmlLines [][]byte
	var writerExist bool

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		golines := Scan(line, &htmlLines, &writerExist)

		if len(golines) > 0 {
			// fmt.Println("===")
			// fmt.Println(string(bytes.Join(golines, []byte("\n"))))

			for _, l := range golines {
				// fmt.Fprintln(w, l)
				w.Write(l)
				w.Write([]byte{'\n'})
				// fmt.Println(">>>", string(l))

			}
			// lines = append(lines, golines...)

			// fmt.Println("\n\n=====")
			// fmt.Println(string(bytes.Join(lines, []byte("\n"))))
		}
	}

	// log.Println(len(htmlLines))

	return nil
}

var (
	reComment = regexp.MustCompile("^[ 	]*//")
	reGoLine  = regexp.MustCompile("^[ 	]*\\|\\|.+")

// 	reInlineTernary = regexp.MustCompile("{{\\?.+?}}")
// 	reInlineGoPrint = regexp.MustCompile("{{=.+?}}")
// 	reInlineSplit   = regexp.MustCompile("(.*?){{(.+?)}}(.*)")
)

//Scan is line parser
func Scan(line []byte, htmlLines *[][]byte, writerExist *bool) (lines [][]byte) {
	//ignore empty line
	if len(bytes.Trim(line, " 	")) == 0 {
		return
	}

	//ignore comment line
	if reComment.Match(line) {
		return
	}

	//go code
	if reGoLine.Match(line) {

		if len(*htmlLines) != 0 {
			lines = append(lines, printHTML(*htmlLines...))
			*htmlLines = nil //[][]byte{}
		}

		lines = append(lines, parseGoLine(line, writerExist)...)
		return
	}

	for {
		n0 := bytes.Index(line, []byte("{{"))
		n1 := bytes.Index(line, []byte("}}"))
		// log.Println(string(line), n0, n1)
		if n0 < 0 && n1 < 0 || n0 >= n1 {
			break
		}

		// log.Println(string(line))

		if len(*htmlLines) != 0 {
			lines = append(lines, printHTML(*htmlLines...))
			*htmlLines = [][]byte{}
		}

		// log.Println(string(line))
		// log.Println(len(line), n0, n1)

		prefixHTML := line[:n0]
		operator := line[n0+2 : n0+3][0]
		gocode := line[n0+3 : n1]
		line = line[n1+2:]

		switch operator {
		case '?':
			lines = append(lines, printHTML(prefixHTML))
			lines = append(lines, parseTernary(gocode)...)
		case '=':
			var suffixHTML []byte
			if !bytes.Contains(line, []byte("{{")) && !bytes.Contains(line, []byte("}}")) {
				suffixHTML = line
				line = nil
			}
			lines = append(lines, printGocode(prefixHTML, gocode, suffixHTML))
		default:
			log.Fatalln("unknown inline operator '%s' in line '%s'", string(operator), string(line))
		}
	}

	if len(line) > 0 {
		*htmlLines = append(*htmlLines, line)
		// lines = append(lines, printHTML(line))
	}

	return
}

func addPackageLine(filename string) []byte {
	filename, _ = filepath.Abs(filename)
	packagename := filepath.Base(filepath.Dir(filename))
	return []byte(fmt.Sprintf("package %s\n", packagename))
}

func addFuncBegin(line []byte) (lines [][]byte, writerExist bool) {
	lines = append(lines, []byte{})
	line = bytes.Replace(line, []byte("template"), []byte("func"), 1)
	if bytes.Contains(line, []byte("w io.Writer")) {
		writerExist = true
		lines = append(lines, append(line, []byte(" {")...))
	} else {
		line = regexp.MustCompile("\\) *\\{* *$").ReplaceAll(line, []byte(") []byte {"))
		writer := []byte("w := new(bytes.Buffer)")
		lines = append(lines, line)
		lines = append(lines, writer)
	}

	return
}

func addFuncEnd(line []byte, writerExist bool) [][]byte {
	if writerExist {
		return [][]byte{[]byte{'}'}}
	}
	return [][]byte{[]byte("return w.Bytes()"), []byte("}")}
}

func addFuncHandler(line []byte) (lines [][]byte) {
	f := regexp.MustCompile("=(.*)").FindAllSubmatch(line, -1)
	// log.Println(f, line)

	if len(f) == 0 || len(f[0]) != 2 {
		log.Fatal(fmt.Errorf("called function not found check your template at string %s", string(line)))
	}

	return [][]byte{[]byte(fmt.Sprintf("w.Write(%s)", f[0][1]))}
}

func parseGoLine(line []byte, writerExist *bool) (lines [][]byte) {
	line = bytes.Trim(line, " 	||")

	var morelines [][]byte
	// var writerExist bool
	// line = line[2:]
	// log.Println(string(line))
	switch {
	case line[0] == '=': //print other template
		morelines = addFuncHandler(line)
	case regexp.MustCompile("^ *template ").Match(line):
		morelines, *writerExist = addFuncBegin(line)
	case regexp.MustCompile("^ *end *$").Match(line):
		morelines = addFuncEnd(line, *writerExist)
	default:
		morelines = append(morelines, line)
	}

	// lines = append(lines, printHTML(*htmlLines...))
	lines = append(lines, morelines...)

	// validateFragment(line, lines...)
	// *htmlLines = [][]byte{}
	return
}

func parseTernary(gocode []byte) (lines [][]byte) {
	nq := bytes.Index(gocode, []byte("?"))
	nc := bytes.Index(gocode, []byte(":"))
	if nq <= 0 {
		log.Fatalln("failed parse ternary operator", string(gocode))
	}

	_if := fmt.Sprintf("if %s {", string(gocode[:nq]))
	var _then string
	if nc == -1 {
		_then = fmt.Sprintf("	fmt.Fprintf(w, \"%%v\", %s)", gocode[nq+1:])
	} else {
		_then = fmt.Sprintf("	fmt.Fprintf(w, \"%%v\", %s)", gocode[nq+1:nc])
	}

	lines = append(lines, []byte(_if))
	lines = append(lines, []byte(_then))
	lines = append(lines, []byte("}"))

	if nc > 0 {
		lines[2] = []byte("} else {")
		_else := fmt.Sprintf("	fmt.Fprintf(w, \"%%v\", %s)", gocode[nc+1:])
		lines = append(lines, []byte(_else))
		lines = append(lines, []byte("}"))
	}

	validateFragment(gocode, lines...)

	return
}

func printGocode(prefixHTML, gocode, suffixHTML []byte) []byte {
	s0 := string(prefixHTML)
	s1 := string(suffixHTML)

	// s0 := strings.Trim(strconv.Quote(string(prefixHTML)), "\"")
	// s1 := strings.Trim(strconv.Quote(string(suffixHTML)), "\"")
	// log.Println(s0, s1)
	s := []byte(fmt.Sprintf("fmt.Fprintf(w, `%s%%v%s`, %s)", s0, s1, gocode))

	validateFragment(gocode, s)
	return s
}

//PrintHTML return write command for html code
func printHTML(htmlLines ...[]byte) []byte {
	if len(htmlLines) == 0 {
		return []byte{}
	}

	htmlCode := bytes.TrimLeft(bytes.Join(htmlLines, []byte("\n")), "\r\n")
	if len(htmlCode) == 0 {
		return []byte{}
	}
	// s := strconv.Quote()
	s := fmt.Sprintf("w.Write([]byte(`%s`))", string(htmlCode))

	validateFragment(htmlCode, []byte(s))

	return []byte(s)
}

func validateFragment(original []byte, lines ...[]byte) {
	opt := &imports.Options{FormatOnly: true, AllErrors: true, Fragment: true}
	_, err := imports.Process("fragment", bytes.Join(lines, []byte("\n")), opt)
	if err != nil {
		fmt.Printf("error in fragment: '%s'\n", string(original))
		fmt.Println(string(bytes.Join(lines, []byte("\n"))))
		log.Fatalln(err)
	}
}
