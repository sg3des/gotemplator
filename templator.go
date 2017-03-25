//Templator is  alternative view on html templates for go!
//is the generator go code returned html

package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/imports"
)

var (
	verbose   = flag.Bool("v", false, "verbose mode")
	extension = flag.String("e", ".gtm", "extension of templates")
	dir       string
)

func init() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()

	dir = flag.Arg(0)
	if dir == "" {
		dir, _ = os.Getwd()
	} else {
		dir, _ = filepath.Abs(dir)
	}

	if _, err := os.Stat(dir); err != nil {
		fmt.Println("directory", dir, "not found")
		os.Exit(1)
	}
}

//main is main
func main() {
	_, err := Generate(dir)
	if err != nil {
		log.Fatal(err)
	}
}

//Generate find all gtm templates in directory, generate go code and save it
func Generate(dir string) (gtms []string, err error) {

	//get all files. with extension .gtm in dir
	gtms, err = filepath.Glob(path.Join(dir, "*"+*extension))
	if err != nil {
		return
	}

	if *verbose {
		fmt.Println("find", len(gtms), "template files")
	}

	//generate it
	for _, gtm := range gtms {
		if *verbose {
			fmt.Println(gtm)
		}
		var filedata []byte

		//parse gtm
		filedata, err = Parse(gtm)
		if err != nil {
			return gtms, fmt.Errorf("failed parse file %s, reason: %s", gtm, err)
		}

		if *verbose {
			scanner := bufio.NewScanner(bytes.NewReader(filedata))
			var i int
			for scanner.Scan() {
				i++
				fmt.Println(i, scanner.Text())
			}
			// fmt.Println(string(filedata))
		}

		// save
		filename := strings.TrimSuffix(gtm, filepath.Ext(gtm)) + ".go"
		// filename := regexp.MustCompile(*extension+"$").ReplaceAllString(gtm, ".go")

		filedata, err = imports.Process(filename, filedata, nil)
		if err != nil {
			return gtms, fmt.Errorf("failed execute `goimport`, error: %s", err)
		}

		if err = ioutil.WriteFile(filename, filedata, 0755); err != nil {
			return
		}
	}
	return
}

//Parse is parser for .gtm file and return go code
func Parse(filename string) ([]byte, error) {

	f, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	if err != nil {
		return []byte{}, err
	}

	var newtemplate []string

	newtemplate = append(newtemplate, addPackageLine(filename))

	scanner := bufio.NewScanner(f)
	var line, prevlines string
	for scanner.Scan() {
		line = scanner.Text()
		if len(line) > 0 {
			line, prevlines = Scan(line, prevlines)
			if len(line) > 0 {
				newtemplate = append(newtemplate, line)
			}
		}
	}

	return []byte(strings.Join(newtemplate, "\n")), nil
}

func addPackageLine(filename string) string {
	packagename := filepath.Base(filepath.Dir(filename))
	return fmt.Sprintf("package %s", packagename)
}

func addWriter(line string) string {
	line = strings.Replace(line, "template", "func", 1)
	line = regexp.MustCompile("\\)\\s*\\{*\\s*$").ReplaceAllString(line, ") []byte {")
	writer := "\n_W := new(bytes.Buffer);\n"

	return line + writer
}

func addReturn(line string) string {
	line = "|| return _W.Bytes()\n}"
	return line
}

func addFuncHandler(line string) string {
	f := regexp.MustCompile("\\|\\|\\=(.*)").FindAllStringSubmatch(line, -1)

	if len(f[0]) != 2 {
		log.Fatal(errors.New("called function not found check your template at string" + line))
	}
	return fmt.Sprintf("|| _W.Write(%s)", f[0][1])
}

// var prevlines string

//Scan is line parser
func Scan(line string, prevlines string) (string, string) {
	//exclude comments
	if regexp.MustCompile("^[ 	]*//").MatchString(line) {
		return "", ""
	}

	if regexp.MustCompile("^[ 	]*\\|\\|").MatchString(line) {

		if regexp.MustCompile("\\|\\|=").MatchString(line) {
			line = addFuncHandler(line)
		}

		if regexp.MustCompile("^\\|\\| *template ").MatchString(line) {
			line = addWriter(line)
		}

		if regexp.MustCompile("^\\|\\| *end *$").MatchString(line) {
			line = addReturn(line)
		}

		return Print(prevlines) + strings.Trim(line, " 	|"), ""
	}

	if regexp.MustCompile("{{.*?}}").MatchString(line) {
		// var ret []string
		aline := regexp.MustCompile("(.*?)({{.*?}})(.*)").FindAllStringSubmatch(line, -1)
		// log.Printf("%#v", aline[0])

		ret := Print(prevlines) + Print(aline[0][1]) + GoPrint(aline[0][2])
		add, prevlines := Scan(aline[0][3], "")

		ret += add + Print(prevlines)

		return ret, ""
	}

	prevlines += strings.TrimLeft(line, " \t\r\n")

	return "", prevlines
}

//Print return write command for html code
func Print(str string) string {
	str = strings.TrimLeft(str, " \t\r\n")
	if str == "" {
		return ""
	}

	return fmt.Sprintf(`_W.WriteString(%s);`, strconv.Quote(str))
}

//GoPrint return go code
func GoPrint(str string) string {
	if regexp.MustCompile("{{=").MatchString(str) {
		val := strings.Trim(str, "{}=")
		return `fmt.Fprintf(_W, "%v", ` + val + `);` + "\n"
	}

	if regexp.MustCompile("{{?").MatchString(str) {
		matches := regexp.MustCompile("{{\\?(.*?)\\?(.*?)(:.*)}}").FindAllStringSubmatch(str, -1)
		if len(matches) == 0 || len(matches[0]) != 4 {
			log.Fatalln("failed parse ternary operator", str)
		}

		_if := matches[0][1]
		_then := matches[0][2]
		_else := strings.TrimLeft(matches[0][3], `:`)

		condition := fmt.Sprintf("\nif %s {\n fmt.Fprintf(_W,\"%%v\",%s); \n}", _if, _then)
		if _else != "" {
			condition += fmt.Sprintf(" else { fmt.Fprintf(_W,\"%%v\",%s); }\n", _else)
		}
		return condition
		// log.Println(condition, then, or)
	}

	return "\n" + strings.Trim(str, "{}")
}
