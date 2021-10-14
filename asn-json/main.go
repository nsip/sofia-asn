package main

import (
	"log"
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

	data, err := os.ReadFile(outpath)
	if err != nil {
		log.Fatalln(err)
	}

	mIdBlock, _ := getIdBlock(string(data))

	inpath4exp := outpath
	outexp := childrenRepl(inpath4exp, mIdBlock)
	// os.WriteFile("./out/asnexp.json", []byte(outexp), os.ModePerm)

	rootWholeBlock := getRootWholeObject(outexp)
	os.WriteFile("./out/asnroot.json", []byte(rootWholeBlock), os.ModePerm)
}
