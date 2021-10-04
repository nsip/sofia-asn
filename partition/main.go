package main

import (
	"fmt"
	"os"
)

func main() {
	outdir := "out"
	os.MkdirAll(fmt.Sprintf("./%s/", outdir), os.ModePerm)

	data, err := os.ReadFile("../data/tree.pretty.json")
	if err != nil {
		panic(err)
	}
	js := string(data)

	la(js, outdir)
	ccp(js, outdir)
	gc(js, outdir)

	//////////////////////////////////////////////////////////////

	data, err = os.ReadFile("../data/node.pretty.json")
	if err != nil {
		panic(err)
	}

	bytesMeta, err := os.ReadFile("../data/metadata.pretty.json")
	if err != nil {
		panic(err)
	}
	mMeta := parseMeta(string(bytesMeta))
	nodeProcess(data, "http://rdf.curriculum.edu.au/202110/", mMeta, outdir)
}
