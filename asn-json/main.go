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

	inpath4exp := outpath
	path4exp := "./out/asnexp.json"

	outexp := childrenRepl(inpath4exp, mIdBlock)
	os.WriteFile(path4exp, []byte(outexp), os.ModePerm)

	/////

	// path4final := "./out/asnfinal.json"
	// outfinal := rmSingleLeaf(outexp)
	// os.WriteFile(path4final, []byte(outfinal), os.ModePerm)
}
