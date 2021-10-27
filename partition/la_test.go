package main

import (
	"os"
	"testing"
)

func TestConnFieldMapping(t *testing.T) {
	data, _ := os.ReadFile("./out/la-English.json")

	bytesMeta, err := os.ReadFile("../data/metadata.pretty.json")
	if err != nil {
		panic(err)
	}
	mMeta := parseMeta(string(bytesMeta))
	// nodeProcess(data, "http://rdf.curriculum.edu.au/202110/", mMeta, outdir)

	js := ConnFieldMapping(string(data), "http://rdf.curriculum.edu.au/202110/", mMeta)
	os.WriteFile("./temp.json", []byte(js), os.ModePerm)
}
