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

	fileContent := ccp(js, outdir)
	err = os.WriteFile(fmt.Sprintf("./%s/ccp-%s.json", outdir, "Cross-curriculum Priorities"), []byte(fileContent), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	for gc, fileContent := range gc(js) {
		err = os.WriteFile(fmt.Sprintf("./%s/gc-%s.json", outdir, gc), []byte(fileContent), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	for la, fileContent := range la(js) {
		err := os.WriteFile(fmt.Sprintf("./%s/la-%s.json", outdir, la), []byte(fileContent), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

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
