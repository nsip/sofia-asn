package main

import (
	"fmt"
	"os"
)

func main() {
	outdir := "out"
	os.MkdirAll(fmt.Sprintf("./%s/", outdir), os.ModePerm)

	bytes, err := os.ReadFile("../data/tree.pretty.json")
	if err != nil {
		panic(err)
	}
	js := string(bytes)

	la(js, outdir)
	ccp(js, outdir)
	gc(js, outdir)

	//////////////////////////////////////////////////////////////

	bytes, err = os.ReadFile("../data/node.pretty.json")
	if err != nil {
		panic(err)
	}

	bytesMeta, err := os.ReadFile("../data/metadata.pretty.json")
	if err != nil {
		panic(err)
	}
	mMeta := parseMeta(string(bytesMeta))
	nodeProcess(bytes, "http://rdf.curriculum.edu.au/202110/", mMeta, outdir)
}
