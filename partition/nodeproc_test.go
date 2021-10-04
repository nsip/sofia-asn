package main

import (
	"fmt"
	"os"
	"testing"
)

func TestParseMeta(t *testing.T) {
	bytes, err := os.ReadFile("../data/metadata.pretty.json")
	if err != nil {
		panic(err)
	}
	mMeta := parseMeta(string(bytes))
	for k, v := range mMeta {
		fmt.Println(k, v)
	}
}

func TestNodeProcess(t *testing.T) {
	bytes, err := os.ReadFile("../data/node.pretty.json")
	if err != nil {
		panic(err)
	}

	bytesMeta, err := os.ReadFile("../data/metadata.pretty.json")
	if err != nil {
		panic(err)
	}
	mMeta := parseMeta(string(bytesMeta))

	nodeProcess(bytes, "http://rdf.curriculum.edu.au/202110/", mMeta, "out")
}
