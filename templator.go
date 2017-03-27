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

	newtemplate := []string{addPackageLine(filename)}

	scanner := bufio.NewScanner(f)

	// var line string
	var htmlLines []string

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			golines := Scan(line, &htmlLines)
			if len(line) > 0 {
				newtemplate = append(newtemplate, golines...)
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
	writer := "\n	_W := new(bytes.Buffer);\n"

	return line + writer
}

func addReturn(line string) string {
	line = "return _W.Bytes()\n}"
	return line
}

func addFuncHandler(line string) string {
	f := regexp.MustCompile("=(.*)").FindAllStringSubmatch(line, -1)
	// log.Println(f, line)

	if len(f) == 0 || len(f[0]) != 2 {
		log.Fatal(errors.New("called function not found check your template at string" + line))
	}
	return fmt.Sprintf("_W.Write(%s)", f[0][1])
}

//Scan is line parser
func Scan(line string, htmlLines *[]string) []string {
	line = strings.Trim(line, " 	")

	if line == "" {
		return []string{}
	}

	//ignore comments
	if len(line) >= 2 && line[:2] == "//" {
		return []string{}
	}

	//go code
	if len(line) >= 2 && line[:2] == "||" {

		line = line[2:]
		switch {
		case line[0] == '=': //print other template
			line = addFuncHandler(line)
		case regexp.MustCompile("^ *template").MatchString(line):
			line = addWriter(line)
		case regexp.MustCompile("^ *end").MatchString(line):
			line = addReturn(line)
		}

		doneLines := []string{PrintHTML(*htmlLines...), line}
		*htmlLines = []string{}
		return doneLines
	}

	//inline go code
	if regexp.MustCompile("{{.*?}}").MatchString(line) {

		aline := regexp.MustCompile("(.*?)({{.*?}})(.*)").FindAllStringSubmatch(line, -1)

		var doneLines []string
		if len(*htmlLines) > 0 {
			doneLines = append(doneLines, PrintHTML(*htmlLines...))
			*htmlLines = []string{}
		}
		doneLines = append(doneLines, PrintHTML(aline[0][1]))
		doneLines = append(doneLines, GoPrint(aline[0][2]))

		add := Scan(aline[0][3], htmlLines)
		if len(add) != 0 {
			doneLines = append(doneLines, add...) //append(add, PrintHTML(*htmlLines...))...)
		}
		if len(*htmlLines) > 0 {
			doneLines = append(doneLines, PrintHTML(*htmlLines...))
			*htmlLines = []string{}
		}

		return doneLines
	}

	//html
	*htmlLines = append(*htmlLines, strings.TrimLeft(line, " \t\r\n"))

	return []string{}
}

//Print return write command for html code
func PrintHTML(htmlLines ...string) string {
	if len(htmlLines) == 0 {
		return ""
	}

	htmlCode := strings.TrimLeft(strings.Join(htmlLines, "\n"), " \t\r\n")
	if htmlCode == "" {
		return ""
	}

	return fmt.Sprintf(`_W.WriteString(%s)`, strconv.Quote(htmlCode))
}

//GoPrint return go code
func GoPrint(str string) string {

	if len(str) > 6 {
		//print variable
		if str[:3] == "{{=" {
			val := strings.Trim(str, "{}=")
			return `fmt.Fprintf(_W, "%v", ` + val + `)`
		}

		//ternary operator
		if str[:3] == "{{?" {
			matches := regexp.MustCompile("{{\\?(.*?)\\?(.*?)(:.*)?}}").FindAllStringSubmatch(str, -1)
			if len(matches) == 0 || len(matches[0]) != 4 {
				log.Fatalln("failed parse ternary operator", str)
			}

			_if := matches[0][1]
			_then := matches[0][2]
			_else := strings.TrimLeft(matches[0][3], `:`)

			condition := []string{
				fmt.Sprintf("if %s {", _if),
				fmt.Sprintf("fmt.Fprintf(_W, \"%%v\", %s)", _then),
				"}",
			}

			if _else != "" {
				conditionElse := []string{
					"} else {",
					fmt.Sprintf("fmt.Fprintf(_W,\"%%v\",%s)", _else),
					"}",
				}
				condition = append(condition[:2], conditionElse...)
			}

			return strings.Join(condition, "\n")
		}
	}

	return strings.Trim(str, "{}")
}
