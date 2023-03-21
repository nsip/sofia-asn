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
	// nodeProcess(data, uri4id, mMeta, outDir)

	js := ConnFieldMapping(string(data), uri4id, mMeta)
	os.WriteFile("./temp.json", []byte(js), os.ModePerm)
}
