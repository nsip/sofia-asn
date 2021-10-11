package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/digisan/gotk/filedir"
)

func main() {

	outdir := "./out/"
	outfile := "asn.json"
	os.MkdirAll(outdir, os.ModePerm)
	outpath := filepath.Join(outdir, outfile)

	if !filedir.FileExists(outpath) {
		data, err := os.ReadFile("../partition/out/node-meta.json")
		if err != nil {
			panic(err)
		}
		nodeProcess(data, outdir, outfile, "../data/tree.pretty.json", "http://rdf.curriculum.edu.au/202110/")
	}

	inpath4exp := outpath
	outpath4exp := "./out/asnexp.json"

	replaced := childrenRepl(inpath4exp, outpath4exp)
	fmt.Println(replaced)
}
