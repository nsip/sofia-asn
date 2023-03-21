package main

import (
	"fmt"
	"os"
	"testing"
)

func TestParseMeta(t *testing.T) {
	data, err := os.ReadFile("../data/Sofia-API-Meta-Data-09062022.json")
	if err != nil {
		panic(err)
	}
	mMeta := parseMeta(string(data))
	for k, v := range mMeta {
		fmt.Println(k, v)
	}
}

func TestNodeProcess(t *testing.T) {
	data, err := os.ReadFile("../data/Sofia-API-Nodes-Data-09062022.json")
	if err != nil {
		panic(err)
	}

	dataMeta, err := os.ReadFile("../data/Sofia-API-Meta-Data-09062022.json")
	if err != nil {
		panic(err)
	}
	mMeta := parseMeta(string(dataMeta))

	nodeProcess(data, uri4id, mMeta, "out")
}
