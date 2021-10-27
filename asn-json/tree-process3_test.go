package main

import (
	"log"
	"os"
	"testing"
)

func TestTreeProc3(t *testing.T) {
	os.MkdirAll("./out/", os.ModePerm)

	file := "la-English.json"
	data, err := os.ReadFile(`../partition/out/` + file)
	if err != nil {
		log.Fatalln(err)
	}
	s := treeProc3(data)
	os.WriteFile("./out/testout.json", []byte(s), os.ModePerm)
}
