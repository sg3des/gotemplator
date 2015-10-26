//go:install
package main

import (
	"os"

	"github.com/sg3des/templator"
)

func main() {
	var globpath = ""
	if len(os.Args) >= 2 && len(os.Args[1]) > 0 {
		globpath = os.Args[1]
	}

	templator.Generate(globpath, true)
}
