package main

import (
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
		nodeProc(data, outdir, outfile, "../data/tree.pretty.json", "http://rdf.curriculum.edu.au/202110/")
	}

	/////

	mIdBlock, _ := getIdBlock(outpath)

	/////

	inpath4exp := outpath
	outpath4exp := "./out/asnexp.json"

	N := 0
AGAIN:
	repl := childrenRepl(inpath4exp, outpath4exp, mIdBlock)
	if repl && N < 100 {
		inpath4exp = outpath4exp
		outpath4exp = outpath4exp[:len(outpath4exp)-5] + "1.json"
		N++
		goto AGAIN
	}
}
