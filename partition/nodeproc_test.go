package main

import (
	"fmt"
	"os"
	"testing"
)

func TestParseMeta(t *testing.T) {
	data, err := os.ReadFile("../data/metadata.pretty.json")
	if err != nil {
		panic(err)
	}
	mMeta := parseMeta(string(data))
	for k, v := range mMeta {
		fmt.Println(k, v)
	}
}

func TestNodeProcess(t *testing.T) {
	data, err := os.ReadFile("../data/node.pretty.json")
	if err != nil {
		panic(err)
	}

	dataMeta, err := os.ReadFile("../data/metadata.pretty.json")
	if err != nil {
		panic(err)
	}
	mMeta := parseMeta(string(dataMeta))

	nodeProcess(data, uri4id, mMeta, "out")
}
