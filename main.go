package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/thisguycodes/copy/reflink"
)

func init() {
	flag.Parse()
}

func main() {
	fs := afero.NewOsFs()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: copy <from> <to>")
		os.Exit(1)
	}

	err := reflink.ReflinkOrCopyAfero(fs, "big.file", "alsobig.file")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
