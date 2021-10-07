package main

import (
	"fmt"
	"os"
)

func main() {

	outdir := "out"
	os.MkdirAll(fmt.Sprintf("./%s/", outdir), os.ModePerm)

	data, err := os.ReadFile("../partition/out/node-meta.json")
	if err != nil {
		panic(err)
	}

	nodeProcess(data, "./out/", "../data/tree.pretty.json", "http://rdf.curriculum.edu.au/202110/")
}
