//Templator is  alternative view on html templates for go!
//is the generator go code returned html

//go:install
package main

import (
	"bufio"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	var globpath = ""
	if len(os.Args) >= 2 && len(os.Args[1]) > 0 {
		globpath = os.Args[1]
	}

	Generate(globpath, true)
}

//Generate find all gtm templates in directory, generate go code and save it
func Generate(dir string, verbose ...bool) {
	//get all files. with extension .gtm in dir
	gtms, err := filepath.Glob(path.Join(dir, "*.gtm"))
	if err != nil {
		panic(err)
	}

	//generate it
	for _, gtm := range gtms {
		if len(verbose) > 0 && verbose[0] {
			fmt.Println(gtm)
		}

		//parse gtm
		filedata, err := Parse(gtm)
		if err != nil {
			panic(err)
		}

		//format code
		filedata, err = format.Source(filedata)
		if err != nil {
			panic(err)
		}

		//save
		filename := regexp.MustCompile("gtm$").ReplaceAllString(gtm, "go")
		if err := ioutil.WriteFile(filename, filedata, 0755); err != nil {
			panic(err)
		}
	}
}

//Parse is parser for .gtm file and return go code
func Parse(filename string) ([]byte, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	if err != nil {
		return []byte{}, err
	}

	var newtemplate []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		newtemplate = append(newtemplate, Scan(line)...)
	}

	return []byte(strings.Join(newtemplate, "\n")), nil
}

//Scan is line parser
func Scan(line string) []string {
	if regexp.MustCompile("^\\|\\|").MatchString(line) {
		return []string{strings.Trim(line, "| ")}
	}

	if regexp.MustCompile("{{.*?}}").MatchString(line) {
		var ret []string
		aline := regexp.MustCompile("(.*?)({{.*?}})(.*)").FindAllStringSubmatch(line, -1)

		ret = append(ret, Print(aline[0][1]))   // print html before go code {{
		ret = append(ret, GoPrint(aline[0][2])) // print go code in {{}}
		ret = append(ret, Scan(aline[0][3])...) // parse string, what is left, after }}

		return ret
	} else {
		return []string{Print(line)}
	}
	return []string{""}
}

//Print return write command for html code
func Print(str string) string {
	return fmt.Sprintf(`w.Write([]byte("%s"))`, str)
}

//GoPrint return go code
func GoPrint(str string) string {
	if regexp.MustCompile("{{=").MatchString(str) {
		val := strings.Trim(str, "{}=")
		return `fmt.Fprintf(w, "%s", ` + val + `)`
	}
	return strings.Trim(str, "{}")
}
